PACKAGE=jsonTail

default: build

build: main.go
	go build \
	    -tags release \
        -o bin/$(PACKAGE) \
        main.go
