.PHONY: test coverage build pre-commit dependencies install-golangci-lint install-pre-commit

test:
	go test ./... -shuffle on -race -coverprofile=coverage.out

coverage: test
	go tool cover -html=coverage.out

build:
	goreleaser build --single-target --snapshot --rm-dist

pre-commit: dependencies
	@pre-commit install

dependencies: install-golangci-lint install-pre-commit

install-golangci-lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2; \
	}

install-pre-commit:
	@command -v asdf >/dev/null 2>&1 && asdf plugin-add pre-commit || true
	@command -v asdf >/dev/null 2>&1 && asdf install pre-commit || pip install pre-commit
