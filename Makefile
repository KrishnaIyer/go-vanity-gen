.PHONY: build

GVG_VERSION=v0.0.1
GVG_GIT_COMMIT=$(shell git rev-parse --short HEAD)
GVG_DATE=$(shell date)
GVG_PACKAGE="go.krishnaiyer.dev/go-vanity-gen"


init:
	@echo "Initialise repository..."
	@mkdir -p gen

test:
	go test ./... -cover

build.local:
	go build \
	-ldflags="-X '${GVG_PACKAGE}/cmd.version=${GVG_VERSION}' \
	-X '${GVG_PACKAGE}/cmd.gitCommit=${GVG_GIT_COMMIT}' \
	-X '${GVG_PACKAGE}/cmd.buildDate=${GVG_DATE}'" main.go

build.dist:
	goreleaser --snapshot --skip-publish --rm-dist

clean:
	@rm -rf dist
	@rm -rf gen
	@mkdir -p gen
