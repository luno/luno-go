.PHONY: install-dev-deps

install-dev-deps:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/go-critic/go-critic/cmd/gocritic@latest
