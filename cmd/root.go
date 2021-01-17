// Copyright © 2020 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.krishnaiyer.dev/go-vanity-gen/pkg/generator"
	cfg "go.krishnaiyer.dev/godry/config"
)

// Template is the input template.
type Template struct {
	Index   string `name:"index" short:"i" description:"path to html template for the index"`
	Project string `name:"project" short:"p" description:"path to html template for projects"`
}

// Config represents the configuration
type Config struct {
	VanityFile string   `name:"file" short:"f" description:"file containing vanity redirection paths (yml)"`
	Template   Template `name:"template"`
	OutPath    string   `name:"out-path" short:"o" description:"directory where output files are generated. Default is ./gen"`
	Debug      bool     `name:"debug" short:"d" description:"print detailed logs for errors"`
}

type kv map[string][]byte

var (
	flags = pflag.NewFlagSet("go-vanity", pflag.ExitOnError)

	config = new(Config)

	manager *cfg.Manager

	addressRegex = regexp.MustCompile(`^([a-z-.0-9]+)(:[0-9]+)?$`)

	errTemplateNotDefined = fmt.Errorf("Template not defined")

	// Root is the root command.
	Root = &cobra.Command{
		Use:           "go-vanity",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "go-vanity generates vanity assets from templates",
		Long:          `go-vanity generates vanity assets from templates. Templates are usually simple html files that contain links to repositories`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := manager.Unmarshal(config)
			if err != nil {
				panic(err)
			}
			if config.OutPath == "" {
				config.OutPath = "./gen"
			}
			if config.Template.Index == "" || config.Template.Project == "" {
				return errTemplateNotDefined
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			baseCtx := context.Background()
			ctx, cancel := context.WithCancel(baseCtx)
			defer cancel()

			input := kv{
				config.VanityFile:       nil,
				config.Template.Index:   nil,
				config.Template.Project: nil,
			}
			for name := range input {
				raw, err := ioutil.ReadFile(name)
				if err != nil {
					log.Fatal(fmt.Errorf("Failed to read file %s: %v", name, err.Error()))
				}
				input[name] = raw
			}

			gen, err := generator.New(ctx, input[config.VanityFile])
			if err != nil {
				log.Fatal(err.Error())
			}

			index, err := gen.Index(ctx, string(input[config.Template.Index]))
			if err != nil {
				log.Fatal(fmt.Errorf("Failed to generate index :%v", err.Error()))
			}
			indexFile := fmt.Sprintf("%s/index.html", config.OutPath)
			err = ioutil.WriteFile(indexFile, index, 0755)
			if err != nil {
				log.Fatal(fmt.Errorf("Failed to write index at %s :%v", indexFile, err.Error()))
			}

			out, err := gen.Project(ctx, string(input[config.Template.Project]))
			if err != nil {
				log.Fatal(fmt.Errorf("Failed to generate project files :%v", err.Error()))
			}

			for name, project := range out.Items() {
				basePath := fmt.Sprintf("%s%s", config.OutPath, name)
				paths := []string{basePath}
				paths = append(paths, project.PkgNames...)
				for _, path := range paths {
					if path != basePath {
						path = basePath + "/" + path
					}
					err := os.MkdirAll(path, 0755)
					if err != nil {
						log.Fatal(fmt.Errorf("Failed to create folder %s :%v", path, err.Error()))
					}
					err = ioutil.WriteFile(path+"/index.html", project.Content, 0755)
					if err != nil {
						log.Fatal(fmt.Errorf("Failed to create file %s :%v", path+"/index.html", err.Error()))
					}
				}
			}
		},
	}
)

// Execute the root command
func Execute() {
	if err := Root.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	manager = cfg.New("config", "go-vanity")
	manager.InitFlags(*config)
	Root.PersistentFlags().AddFlagSet(manager.Flags())
	Root.AddCommand(manager.VersionCommand(Root))
	Root.AddCommand(manager.ConfigCommand(Root))
}
