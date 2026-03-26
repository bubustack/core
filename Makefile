ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: check-tractatus
check-tractatus: ## Verify a sibling tractatus checkout only when go.mod uses a local replace.
	@if grep -Eq '^replace[[:space:]]+github.com/bubustack/tractatus[[:space:]]+=>[[:space:]]+\.\./tractatus$$' go.mod; then \
		test -d ../tractatus || (echo "expected sibling ../tractatus checkout because go.mod locally replaces github.com/bubustack/tractatus => ../tractatus" && exit 1); \
	fi

.PHONY: fmt
fmt: ## Run go fmt against the repository.
	go fmt ./...

.PHONY: fmt-check
fmt-check: ## Fail if any Go files need formatting.
	@fmt_out=$$(gofmt -l .); \
	if [ -n "$$fmt_out" ]; then \
		echo "The following files need gofmt:"; \
		echo "$$fmt_out"; \
		exit 1; \
	fi

.PHONY: vet
vet: check-tractatus ## Run go vet against the repository.
	go vet ./...

.PHONY: test
test: check-tractatus ## Run unit tests with the race detector.
	go test -race ./...

.PHONY: test-coverage
test-coverage: check-tractatus ## Run tests with a coverage profile.
	go test -coverprofile=coverage.out ./...
	@echo "Coverage profile written to coverage.out"

.PHONY: build
build: check-tractatus ## Build all packages.
	go build ./...

.PHONY: clean
clean: ## Remove local build artifacts.
	rm -f coverage.out coverage.html
	go clean ./...

.PHONY: tidy
tidy: ## Run go mod tidy.
	go mod tidy

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

GOLANGCI_LINT = $(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.11.4

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.

$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))
	@test -f .custom-gcl.yml && { \
		echo "Building custom golangci-lint with plugins..." && \
		$(GOLANGCI_LINT) custom --destination $(LOCALBIN) --name golangci-lint-custom && \
		mv -f $(LOCALBIN)/golangci-lint-custom $(GOLANGCI_LINT); \
	} || true

.PHONY: lint
lint: check-tractatus golangci-lint ## Run golangci-lint.
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: check-tractatus golangci-lint ## Run golangci-lint with fixes enabled.
	$(GOLANGCI_LINT) run --fix

.PHONY: lint-config
lint-config: golangci-lint ## Verify golangci-lint configuration.
	$(GOLANGCI_LINT) config verify

define go-install-tool
@[ -f "$(1)-$(3)" ] && [ "$$(readlink -- "$(1)" 2>/dev/null)" = "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $$(realpath $(1)-$(3)) $(1)
endef
