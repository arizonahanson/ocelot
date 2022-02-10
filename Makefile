NAME=ocelot
PKGPATH=github.com/starlight/ocelot
GITHASH=$(shell git rev-parse --short=12 HEAD)
UTCTIME=$(shell date -u '+%Y%m%d%H%M%S')
LDFLAGS=-X $(PKGPATH)/cmd.version=$(UTCTIME)-$(GITHASH)
GONAME=$(HOME)/go/bin/$(NAME)
GODEPS=main.go cmd/*.go internal/**/*.go pkg/**/*.go Makefile
PIGEON=pigeon
PEGIN=internal/parser/parser.peg
PEGOUT=internal/parser/parser.go

default: install

$(NAME): $(GODEPS)
	@echo "Building..."
	go build -ldflags="$(LDFLAGS)"

$(GONAME): $(GODEPS) $(PEGOUT)
	@echo "Compiling and installing..."
	go install -ldflags="$(LDFLAGS)"

$(PEGOUT): $(PEGIN)
	@echo "Generating parser..."
	$(PIGEON) -o "$(PEGOUT)" "$(PEGIN)"

.PHONY: build
build: $(NAME)

.PHONY: install
install: $(GONAME)
