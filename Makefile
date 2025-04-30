VERSION=$(shell git rev-parse --short HEAD)

.PHONY: default
default: slide

.PHONY: all
all: clean slide

slide:
	@CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w -X github.com/jdel/slide/options.Version=${VERSION}"

.PHONY: clean
clean:
	-@rm -rf slide

REDIST_LICENSES.md:
	@go-licenses report --template redist-licenses.tpl --ignore "github.com/jdel/slide" . > REDIST_LICENSES.md