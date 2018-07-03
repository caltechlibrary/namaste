#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = namaste

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\`  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PKGASSETS = $(shell which pkgassets)

PROJECT_LIST = namaste

OS = $(shell uname)

EXT = 
ifeq ($(OS), Windows)
	EXT = .exe
endif


namaste$(EXT): bin/namaste$(EXT)

cmd/namaste/assets.go:
	pkgassets -o cmd/namaste/assets.go -p main -ext=".md" -strip-prefix="/" -strip-suffix=".md" Examples how-to Help docs/namaste
	git add cmd/namaste/assets.go

bin/namaste$(EXT): namaste.go char_encoding.go cmd/namaste/namaste.go cmd/namaste/assets.go
	go build -o bin/namaste$(EXT) cmd/namaste/namaste.go cmd/namaste/assets.go

build: $(PROJECT_LIST)

install: 
	env GOBIN=$(GOPATH)/bin go install cmd/namaste/namaste.go cmd/namaste/assets.go

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css
	bash mk-website.bash

test: clean bin/namaste$(EXT)
	go test

format:
	gofmt -w namaste.go
	gofmt -w namaste_test.go
	gofmt -w char_encoding.go
	gofmt -w char_encoding_test.go
	gofmt -w cmd/namaste/namaste.go

lint:
	golint namaste.go
	golint namaste_test.go
	golint char_encoding.go
	golint char_encoding_test.go
	golint cmd/namaste/namaste.go

clean: 
	if [ "$(PKGASSETS)" != "" ]; then bash rebuild-assets.bash; fi
	if [ -f index.html ]; then rm *.html; fi
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi
	if [ -d testdata ]; then rm -fR testdata; fi

dist/linux-amd64:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/namaste cmd/namaste/namaste.go cmd/namaste/assets.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/namaste.exe cmd/namaste/namaste.go cmd/namaste/assets.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/namaste cmd/namaste/namaste.go cmd/namaste/assets.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/namaste cmd/namaste/namaste.go cmd/namaste/assets.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

distribute_docs:
	if [ -d dist ]; then rm -fR dist; fi
	mkdir -p dist
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	bash package-versions.bash > dist/package-versions.txt

update_version:
	./update_version.py --yes

release: clean namaste.go char_encoding.go distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

publish:
	bash mk-website.bash
	bash publish.bash

