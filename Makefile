DIRS := . ./storage

clean:
	rm -f ./skelet

build:
	go build -v

test:
	go test -v ${DIRS}

install:
	go install

check:
	go vet
	golint

run: build
	./skelet
