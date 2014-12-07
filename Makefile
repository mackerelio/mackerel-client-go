all: test

test: testdeps
	go test -v ./...

testdeps:
	go get -d -v -t ./...

.PHONY: all test testdeps
