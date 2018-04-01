NAME := kcl
#VERSION := $(shell git describe --tags --abbrev=0)
VERSION := $(shell git rev-list --count HEAD)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGSLDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'

# 必要なツール類をセットアップする
## Setup
setup:
  go get github.com/Masterminds/glide
  go get github.com/golang/lint/golint
  go get golang.org/x/tools/cmd/goimports
  go get github.com/Songmu/make2help/cmd/make2help
  go get github.com/harlow/kinesis-consumer
  go get github.com/harlow/kinesis-consumer/checkpoint/ddb

# テストを実行する
## Run tests
test: deps
  go test $$(glide novendor)

# glideを使って依存パッケージをインストールする
## Install dependencies
deps: setup
  glide install

## Update dependencies
update: setup
  glide update

## Lint
lint: setup
  go vet $$(glide novendor)
  for pkg in $$(glide novendor -x); do \
    golint --set_exit_status $$pkg || exit $$?; \
  done

## Format source codes
fmt: setup
  goimports -w $$(glide nv -x)

## build binaries ex. make bin/myproj
bin/%: cmd/%/main.go deps
  go build -ldflags "$(LDFLAGS)" -o $@ $<

## Show help
help:
  @make2help $(MAKEFILE_LIST)

.PHONY: setup deps update test lint help