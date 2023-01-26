
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: fmt vet
	go test -coverprofile=coverage.out github.com/maxwu/chainer/... -v -ginkgo.v -ginkgo.progress -test.v

.PHONY: coverage
coverage:
	go tool cover -html coverage.out -o coverage.html
