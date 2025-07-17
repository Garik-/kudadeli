GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
GIT_HASH ?= $(shell git rev-parse --short HEAD)

VERSION := $(GIT_BRANCH)-$(GIT_HASH)
GO_LDFLAGS := -ldflags="-w -s -X 'main.Version=$(VERSION)'"

BINDIR := $(CURDIR)/bin
DISTDIR := $(CURDIR)/dist

LINT_VERSION ?= v2.2.2
LINT_BIN := $(BINDIR)/golangci-lint

APP := kudadeli
ARCHS := arm64 arm


.PHONY: tidy lint lint-install fix vuln test run build

tidy:
	@go mod tidy

lint: $(LINT_BIN)
	$(LINT_BIN) run ./...

fix: $(LINT_BIN)
	$(LINT_BIN) run --fix ./...

$(LINT_BIN):
	@echo "Installing golangci-lint $(LINT_VERSION)..."
	@mkdir -p $(BINDIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(BINDIR) $(LINT_VERSION)

vuln: tidy
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

test: tidy
	@go test -v -race -count=5 ./...

run: tidy
	@export $(shell xargs < .env) && go run .


build: tidy
	@echo "Building..."
	@cd ../frontend && npm install && npm run build-only
	@CGO_ENABLED=0 go build $(GO_LDFLAGS) -o $(DISTDIR)/$(APP) .
	@for arch in $(ARCHS); do \
		echo "Building for $$arch..."; \
		CGO_ENABLED=0 GOOS=linux GOARCH=$$arch go build $(GO_LDFLAGS) -o $(DISTDIR)/$(APP)-$$arch .; \
	done
	@ls -lh $(DISTDIR)