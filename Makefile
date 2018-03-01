PACKAGE=jsonTail

default: build

build: src/*.go
	go build \
	    -tags release \
        -o bin/$(PACKAGE) \
        src/main.go
