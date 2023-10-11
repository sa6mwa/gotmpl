VERSION = 0.1.0

.PHONY: clean install release

all: bin/gotmpl

clean:
	rm -rf bin
	rm -f gotmpl-$(shell go env GOOS)-$(shell go env GOARCH)-$(VERSION).tar.gz

install: bin/gotmpl
	sudo install bin/gotmpl /usr/local/bin/

release: bin/gotmpl
	sha256sum bin/gotmpl > CHECKSUMS.sha256sum
	tar -czf gotmpl-$(shell go env GOOS)-$(shell go env GOARCH)-$(VERSION).tar.gz --exclude-vcs --owner=0 --group=0 --transform 's|^|gotmpl-$(VERSION)/|' *.go LICENSE README.md Makefile go.mod go.sum CHECKSUMS.* bin/gotmpl


bin:
	mkdir -p bin

bin/gotmpl: bin
	CGO_ENABLED=0 go test -v -cover ./...
	CGO_ENABLED=0 go build -o bin/gotmpl
	strip -s bin/gotmpl
	if which upx > /dev/null ; then upx -9 bin/gotmpl ; fi
