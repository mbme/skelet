clean:
	rm -f ./skelet

build:
	go build -v

test: build
	go test

check:
	go vet
	golint

run: build
	./skelet
