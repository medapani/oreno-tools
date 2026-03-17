WAILS ?= wails

.PHONY: build go-build build-mac-win test cover

build:
	$(WAILS) build -trimpath -ldflags="-s -w"

build-win:
	$(WAILS) build -platform "windows/amd64" -trimpath -ldflags="-s -w" -nopackage

# リント
lint:
	@echo "Linting..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Installing..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# テスト
test:
	go test ./backend/... -v -count=1 -timeout=120s

# カバレッジ (HTML レポートをブラウザで表示)
cover:
	go test ./backend/... -count=1 -timeout=120s -coverprofile=coverage.out
	go tool cover -html=coverage.out
