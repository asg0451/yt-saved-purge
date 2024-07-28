srcs = $(shell find . -name '*.go' | grep -v '.pb.go')

.PHONY: all build gen push-image deploy

all: build

build: $(srcs) gen
	mkdir -p bin
	CGO_ENABLED=1 GOOS=linux CC=`which musl-gcc` go build -ldflags='-linkmode external -s -w -extldflags "-static"' -o bin/ ./...

gen:
	go generate ./...

push-image:
	gh workflow run push-image.yaml
	sleep 5
	gh run watch "$$(gh run list --workflow=push-image.yaml | cut -d'	' -f7 | head -n1)"

deploy: push-image
