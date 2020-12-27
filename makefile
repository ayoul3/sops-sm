BUILD=go build -ldflags="-s -w"
TARGETDIR=""
build:
	$(BUILD) sops-sm.go

run:
	go run sops-sm.go

install:
	go install sops-sm.go

test:
	go test ./...

linux:
	GOOS=linux GOARCH=amd64 $(BUILD) -o sops-sm sops-sm.go

windows:
	GOOS=windows GOARCH=amd64 $(BUILD) -o sops-sm-windows.exe sops-sm.go

darwin:
	GOOS=darwin GOARCH=amd64 $(BUILD) -o sops-sm-darwin sops-sm.go

freebsd:
	GOOS=freebsd GOARCH=amd64 $(BUILD) -o sops-sm-freebsd sops-sm.go

build-all: linux windows darwin freebsd