DIRS := . ./storage
BASE := .
GO_FILES = $(wildcard *.go)

clean:
	rm -f $(BASE)/skelet

build:
	go build -v $(BASE)

test:
	go test -v ${DIRS}

install:
	go install $(BASE)

check:
	go vet $(BASE)
	golint $(BASE)

run: build
	$(BASE)/skelet

serv:
	node $(BASE)/watch.js $(BASE)

.PHONY: clean build test install check run serv
