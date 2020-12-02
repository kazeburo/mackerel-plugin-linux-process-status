VERSION=0.0.8
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
all: mackerel-plugin-linux-process-status

.PHONY: mackerel-plugin-linux-process-status

mackerel-plugin-linux-process-status: linux.go
	go build $(LDFLAGS) -o mackerel-plugin-linux-process-status

linux: linux.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-linux-process-status

fmt:
	go fmt ./...

check:
	go test ./...

clean:
	rm -rf mackerel-plugin-linux-process-status

tag:
	git tag v${VERSION}
	git push origin v${VERSION}
	git push origin master
