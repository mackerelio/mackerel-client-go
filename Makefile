.PHONY: test
test: lint gofmt
	go test -v ./...

.PHONY: testdeps
testdeps:
	go get -d -v -t ./...
	GO111MODULE=off \
	go get golang.org/x/lint/golint \
		golang.org/x/tools/cmd/cover \
		github.com/axw/gocov/gocov \
		github.com/mattn/goveralls

LINT_RET = .golint.txt
.PHONY: lint
lint: testdeps
	go vet .
	rm -f $(LINT_RET)
	golint ./... | tee $(LINT_RET)
	test ! -s $(LINT_RET)

GOFMT_RET = .gofmt.txt
.PHONY: gofmt
gofmt: testdeps
	rm -f $(GOFMT_RET)
	gofmt -s -d *.go | tee $(GOFMT_RET)
	test ! -s $(GOFMT_RET)

.PHONY: cover
cover: testdeps
	goveralls
