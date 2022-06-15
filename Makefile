.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test: fmt
	ginkgo -r -cover