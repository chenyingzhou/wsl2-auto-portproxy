.PHONY:build clean dev
# Go parameters
GO=go

# name
BINARY_NAME=wsl2-tcpproxy

build: mod-tidy
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -o ./$(BINARY_NAME).exe

clean:
	$(GO) clean
	@rm -f ./$(BINARY_NAME).exe

mod-tidy:
	$(GO) mod tidy

dev:
	$(GO) run ./main.go
