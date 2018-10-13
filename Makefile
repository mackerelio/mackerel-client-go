test: lint gofmt
	go test -v ./...

testdeps:
	go get -d -v -t ./...
	go get golang.org/x/lint/golint \
		golang.org/x/tools/cmd/cover \
		github.com/axw/gocov/gocov \
		github.com/mattn/goveralls

LINT_RET = .golint.txt
lint: testdeps
	go tool vet .
	rm -f $(LINT_RET)
	golint ./... | tee $(LINT_RET)
	test ! -s $(LINT_RET)

GOFMT_RET = .gofmt.txt
gofmt: testdeps
	rm -f $(GOFMT_RET)
	gofmt -s -d *.go | tee $(GOFMT_RET)
	test ! -s $(GOFMT_RET)

cover: testdeps
	goveralls

.PHONY: test testdeps lint gofmt cover
