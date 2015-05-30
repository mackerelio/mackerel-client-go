all: test

test: lint
	go test -v ./...

testdeps:
	go get -d -v -t ./...
	go get golang.org/x/tools/cmd/vet
	go get github.com/golang/lint/golint

LINT_RET = .golint.txt
lint: testdeps
	go tool vet .
	rm -f $(LINT_RET)
	golint ./... | tee .golint.txt
	test ! -s $(LINT_RET)

.PHONY: all test testdeps
