.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test: fmt
	ginkgo -r -cover

.PHONY: lint
lint:
	golangci-lint run