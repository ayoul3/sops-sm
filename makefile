build:
	go build sops-sm.go

run:
	go run sops-sm.go

install:
	go install sops-sm.go

test:
	go test ./...

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 main.go

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows-386 main.go

darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 main.go

freebsd:
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd-amd64 main.go