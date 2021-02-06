// Copyright © 2021 Krishna Iyer Easwaran
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

package generator

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"

	"gopkg.in/yaml.v2"
)

// Path is the parsed configuation of vanity paths.
type Path struct {
	path     string
	repo     string
	display  string
	vcs      string
	packages []string
}

type pathConfigSet []Path

// Generator generates vanity assests.
type Generator struct {
	cfg   config
	paths []Path
	host  string
}

// OutItem is a single output item.
type OutItem struct {
	PkgNames []string
	Content  []byte
}

// Out is the raw output from the generator.
type Out struct {
	items map[string]OutItem
}

// Items returns the generated output items.
func (o *Out) Items() map[string]OutItem {
	return o.items
}

// config is the vanity config
type config struct {
	Host  string `yaml:"host,omitempty"`
	Paths map[string]struct {
		Repo     string   `yaml:"repo,omitempty"`
		Display  string   `yaml:"display,omitempty"`
		VCS      string   `yaml:"vcs,omitempty"`
		Packages []string `yaml:"packages,omitempty"`
	} `yaml:"paths,omitempty"`
}

// New parses the provided vanity config and returns a new Generator.
func New(ctx context.Context, vanity []byte) (*Generator, error) {
	var vanityCfg config
	if err := yaml.Unmarshal(vanity, &vanityCfg); err != nil {
		return nil, fmt.Errorf("Could not parse vanity config: %v", err)
	}
	paths := make([]Path, 0)
	for path, e := range vanityCfg.Paths {
		pc := Path{
			path:     strings.TrimSuffix(path, "/"),
			repo:     e.Repo,
			display:  e.Display,
			vcs:      e.VCS,
			packages: e.Packages,
		}
		switch {
		case e.Display != "":
		case strings.HasPrefix(e.Repo, "https://github.com/"):
			pc.display = fmt.Sprintf("%v %v/tree/master{/dir} %v/blob/master{/dir}/{file}#L{line}", e.Repo, e.Repo, e.Repo)
		case strings.HasPrefix(e.Repo, "https://bitbucket.org"):
			pc.display = fmt.Sprintf("%v %v/src/default{/dir} %v/src/default{/dir}/{file}#{file}-{line}", e.Repo, e.Repo, e.Repo)
		}
		switch {
		case e.VCS != "":
			if e.VCS != "bzr" && e.VCS != "git" && e.VCS != "hg" && e.VCS != "svn" {
				return nil, fmt.Errorf("configuration for %v: unknown VCS %s", path, e.VCS)
			}
		case strings.HasPrefix(e.Repo, "https://github.com/"):
			pc.vcs = "git"
		default:
			return nil, fmt.Errorf("configuration for %v: cannot infer VCS from %s", path, e.Repo)
		}
		paths = append(paths, pc)
	}
	return &Generator{
		cfg:   vanityCfg,
		paths: paths,
		host:  vanityCfg.Host,
	}, nil
}

type vanity struct {
	Path string
	Repo string
}

// Index generates the index.html at the root of the assets tree.
func (gen *Generator) Index(ctx context.Context, input string) ([]byte, error) {
	index, err := template.New("index").Parse(input)
	if err != nil {
		return nil, err
	}
	vanityPaths := make([]vanity, len(gen.paths))
	for i, h := range gen.paths {
		vanityPaths[i].Path = gen.host + h.path
		vanityPaths[i].Repo = h.repo
	}
	var buf bytes.Buffer
	if err := index.Execute(&buf, struct {
		Host   string
		Vanity []vanity
	}{
		Host:   gen.host,
		Vanity: vanityPaths,
	},
	); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Project generates the index.html for a project path and returns an Out object.
func (gen *Generator) Project(ctx context.Context, input string) (*Out, error) {
	out := &Out{
		items: make(map[string]OutItem, 0),
	}
	vanity, err := template.New("vanity").Parse(input)
	if err != nil {
		return nil, err
	}
	for _, path := range gen.paths {
		var buf bytes.Buffer
		if err := vanity.Execute(&buf, struct {
			Import  string
			Subpath string
			Repo    string
			Display string
			VCS     string
			Host    string
		}{
			Import:  gen.host + path.path,
			Repo:    path.repo,
			Display: path.display,
			VCS:     path.vcs,
			Host:    gen.host,
		}); err != nil {
			return nil, err
		}
		out.items[path.path] = OutItem{
			PkgNames: path.packages,
			Content:  buf.Bytes(),
		}
	}
	return out, nil
}

// Paths returns the list of paths extracted from the vanity configuration.
func (gen *Generator) Paths() []Path {
	return gen.paths
}
