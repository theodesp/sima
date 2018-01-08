.PHONY: format
format:
	@find . -type f -name "*.go*" -print0 | xargs -0 gofmt -s -w

.PHONY: debs
debs:
	GOPATH=$(GOPATH) go get ./...
	GOPATH=$(GOPATH) go get -u gopkg.in/check.v1
	GOPATH=$(GOPATH) go get -u github.com/OneOfOne/cmap
	GOPATH=$(GOPATH) go get -u github.com/fortytw2/leaktest
	GOPATH=$(GOPATH) go get -u github.com/spaolacci/murmur3
	GOPATH=$(GOPATH) go get -u github.com/deckarep/golang-set

.PHONY: test
test:
	GOPATH=$(GOPATH) go test -race

.PHONY: bench
bench:
	GOPATH=$(GOPATH) go test -bench=. -check.b -benchmem

# Clean junk
.PHONY: clean
clean:
	GOPATH=$(GOPATH) go clean ./...
