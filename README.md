# Go Vanity Generator

Generates simple html pages for Go Vanity redirection.

![Status Checks](https://github.com/krishnaiyer/go-vanity-gen/workflows/buildandtest/badge.svg)


## Installation

```bash
$ go get -u go.krishnaiyer.dev/go-vanity-gen
```

## Options
```
go-vanity generates vanity assets from templates. Templates are usually simple html files that contain links to repositories

Usage:
  go-vanity [flags]
  go-vanity [command]

Available Commands:
  config      Display raw config values
  help        Help about any command
  version     Display version information

Flags:
  -d, --debug                     print detailed logs for errors
  -f, --file string               file containing vanity redirection paths (yml)
  -h, --help                      help for go-vanity
  -o, --out-path string           directory where output files are generated. Default is ./gen
  -i, --template.index string     path to html template for the index
  -p, --template.project string   path to html template for projects

Use "go-vanity [command] --help" for more information about a command.
```

## Usage

```bash
$ go-vanity-gen -i sample/index.tmpl -p sample/project.tmpl -f ./vanity.yml
```

## Development

The following are the prerequisites;
* [Go](https://golang.org/) installed and configured.
* [optional] Custom HTML templates.

1. Clone this repository.
2. Initialize
```bash
$ make init
```
3. Build from source
```bash
$ GVG_PACKAGE=<your-path>/go-vanity-gen GVG_VERSION=<version> make build.local
```
4. Run tests
```bash
$ make test
```
5. Clean up
```bash
$ make clean
```

## Limitations

* Since we're building static assets, each package that need redirection needs an `index.html`. This does result in a lot of duplication. But since each file is very small, the cost of storing these files is much lower than running a server to serve these paths.
* Package paths need to be manually listed. In the root of each of your repositories, use the following command.
```bash
$ go list ./...
```

## License

The contents of this repository are provided "as-is" under the terms of the [Apache 2.0 License](./LICENSE).
