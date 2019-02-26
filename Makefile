BINARY_NAME=bin/hawk.catcher
BINARY_NAME_LINUX=$(BINARY_NAME)-linux
BINARY_NAME_WINDOWS=$(BINARY_NAME)-windows.exe
BINARY_NAME_DARWIN=$(BINARY_NAME)-darwin

export GO111MODULE=on

all: check lint build

build:
	go build -o $(BINARY_NAME) -v ./
check:
lint:
	golint cmd/... lib/... ./
clean:
	go clean
	rm -rf $(BINARY_NAME)
	rm -rf $(BINARY_NAME_LINUX)
	rm -rf $(BINARY_NAME_WINDOWS)
	rm -rf $(BINARY_NAME_DARWIN)

build-all: build-linux build-windows build-darwin

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME_LINUX) -v $(SRC_DIRECTORY)

build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME_WINDOWS) -v $(SRC_DIRECTORY)

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME_DARWIN) -v $(SRC_DIRECTORY)