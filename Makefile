VERSION := $(shell grep "version = " cmd/version.go | awk '{print $$4}' | sed 's/"//g')

PLATFORMS := linux/amd64 darwin/amd64 windows/amd64
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

COMMONENVVAR = GOOS=$(shell uname -s | tr A-Z a-z) GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m)))
BUILDENVVAR = CGO_ENABLED=0

all: build

deps:
	git config --global url."git@github.com:".insteadOf "https://github.com/"
	dep ensure

test:
	$(COMMONENVVAR) $(BUILDENVVAR) go test ./... -v

build:
	$(COMMONENVVAR) $(BUILDENVVAR) go vet ./...
	$(COMMONENVVAR) $(BUILDENVVAR) go build -o nr *.go

tag: # used in client side to trigger a travis deployment job
ifndef VERSION
	$(error VERSION is undefined - run using make tag VERSION=vX.Y.Z)
endif
	git tag $(VERSION)

	# Check to make sure the tag isn't "-dirty".
	if git describe --tags --dirty | grep dirty; \
	then echo current git working tree is "dirty". Make sure you do not have any uncommitted changes ;false; fi

	git push --tags

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) $(BUILDENVVAR) go build -o dist/nr-$(os)-$(arch) *.go

release: $(PLATFORMS)

.PHONY: all deps test build tag release $(PLATFORMS)
