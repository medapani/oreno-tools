WAILS ?= wails

.PHONY: build go-build build-mac-win

build:
	$(WAILS) build -trimpath -ldflags="-s -w"

build-win:
	$(WAILS) build -platform "windows/amd64" -trimpath -ldflags="-s -w" -nopackage

# リント
lint:
	@echo "Linting..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Installing..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...
