PKGPATH=github.com/starlight/ocelot
GITHASH=$(shell git rev-parse --short HEAD)
UTCTIME=$(shell date -u '+%Y%m%d%H%M%S')
LDFLAGS=-X $(PKGPATH)/cmd.version=$(UTCTIME)-$(GITHASH)

GODEPS=main.go cmd/*.go Makefile

.PHONY: default
default: install

.PHONY: build
build: $(GODEPS)
	go build -ldflags="$(LDFLAGS)"

.PHONY: install
install: $(GODEPS)
	go install -ldflags="$(LDFLAGS)"
