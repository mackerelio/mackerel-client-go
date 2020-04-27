.PHONY: test
test: lint gofmt
	go test -v ./...

.PHONY: testdeps
testdeps:
	go install \
		golang.org/x/lint/golint \
		golang.org/x/tools/cmd/cover

.PHONY: lint
lint: testdeps
	golint -set_exit_status ./...

.PHONY: gofmt
gofmt: testdeps
	! gofmt -s -d ./ | grep '^'
