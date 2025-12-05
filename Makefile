.PHONY: build install clean test

build:
	go build -o bin/reqstat .

install:
	go install .

clean:
	rm -rf bin/

test:
	go test -v ./...

run:
	go run . get https://api.github.com/users/octocat

