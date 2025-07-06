#!make
-include .env

.DEFAULT_GOAL := build
build:
	go build


test:
	./scripts/test.sh ./...

#Run code check with all golangci-lint checkers
ALLOWED_BOX_CHARS := "├\|│\|└\|┌\|┐\|┘\|┴\|┬\|┤\|┼"
lint:
	@if LC_ALL=C grep -rn "//.*[^	 !-~]" . --include="*.go" --exclude-dir=vendor | grep -v $(ALLOWED_BOX_CHARS) | grep -q .; then \
		echo "ERROR: Found non-English characters in comments:"; \
		LC_ALL=C grep -rn "//.*[^	 !-~]" . --include="*.go" --exclude-dir=vendor | grep -v $(ALLOWED_BOX_CHARS); \
		exit 1; \
	fi
	GOFUMPT_SPLIT_LONG_LINES=on golangci-lint run