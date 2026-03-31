VERSION=0.1.1

build:
	GOOS=linux GOARCH=386 go build -o bin/msspr-$(VERSION)_linux_386
	GOOS=linux GOARCH=arm64 go build -o bin/msspr-$(VERSION)_linux_arm64
	GOOS=linux GOARCH=386 go build -o bin/msspr-$(VERSION)_linux_386  