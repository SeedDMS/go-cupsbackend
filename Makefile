PKGNAME=printer-driver-seeddms
VERSION=0.0.2

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOGET=$(GOCMD) get

default: build

build: test
	mkdir -p bin
	$(GOBUILD) -v -o bin/cups-seeddms ./cupsbackend/

clean:
	$(GOCLEAN)
	rm -rf bin

test:
	$(GOVET) ./...
	$(GOTEST) -v ./...

run:
	bin/cups-seeddms job-id user title copies options Makefile

stdin:
	cat /tmp/mmk.ps | bin/cups-seeddms job-id user titlestdin copies options

dist: clean
	rm -rf ${PKGNAME}-${VERSION}
	mkdir ${PKGNAME}-${VERSION}
	cp -r cupsbackend *.go *.ppd Makefile README.md go.mod go.sum seeddms-cups.yaml ${PKGNAME}-${VERSION}
	tar czvf ${PKGNAME}-${VERSION}.tar.gz ${PKGNAME}-${VERSION}
	rm -rf ${PKGNAME}-${VERSION}

debian: dist
	mv ${PKGNAME}-${VERSION}.tar.gz ../${PKGNAME}_${VERSION}.orig.tar.gz
	debuild

.PHONY: build test debian
