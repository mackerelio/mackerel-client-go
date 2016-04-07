test: lint
	go test -v ./...

testdeps:
	go get -d -v -t ./...
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/axw/gocov/gocov
	go get github.com/mattn/goveralls

LINT_RET = .golint.txt
lint: testdeps
	go tool vet .
	rm -f $(LINT_RET)
	golint ./... | tee .golint.txt
	test ! -s $(LINT_RET)

cover: testdeps
	goveralls

.PHONY: test testdeps lint cover
