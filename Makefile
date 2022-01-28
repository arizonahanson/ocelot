NAME=ocelot
PKGPATH=github.com/starlight/ocelot
GITHASH=$(shell git rev-parse --short HEAD)
UTCTIME=$(shell date -u '+%Y%m%d%H%M%S')
LDFLAGS=-X $(PKGPATH)/cmd.version=$(UTCTIME)-$(GITHASH)
GONAME=$(HOME)/go/bin/$(NAME)
GODEPS=main.go cmd/*.go Makefile

default: install

$(NAME): $(GODEPS)
	go build -ldflags="$(LDFLAGS)"

$(GONAME): $(GODEPS)
	go install -ldflags="$(LDFLAGS)"

.PHONY: build
build: $(NAME)

.PHONY: install
install: $(GONAME)
