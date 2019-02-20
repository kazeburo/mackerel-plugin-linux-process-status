VERSION=0.0.3
LDFLAGS=-ldflags "-X main.Version=${VERSION}"
all: mackerel-plugin-linux-process-status

.PHONY: mackerel-plugin-linux-process-status

bundle:
	dep ensure

update:
	dep ensure -update

mackerel-plugin-linux-process-status: linux.go
	go build $(LDFLAGS) -o mackerel-plugin-linux-process-status

linux: linux.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-linux-process-status

fmt:
	go fmt ./...

clean:
	rm -rf mackerel-plugin-linux-process-status

tag:
	git tag v${VERSION}
	git push origin v${VERSION}
	git push origin master
	goreleaser --rm-dist
