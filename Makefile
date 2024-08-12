BINARY_NAME=hikup
VERSION=1.0.0
PACKAGE_NAME=$(BINARY_NAME)_$(VERSION)_amd64

all: build

build:
	go build -o $(BINARY_NAME) main.go

package: build
	mkdir -p $(PACKAGE_NAME)/DEBIAN
	mkdir -p $(PACKAGE_NAME)/usr/bin
	mkdir -p $(PACKAGE_NAME)/etc/systemd/system
	cp $(BINARY_NAME) $(PACKAGE_NAME)/usr/bin/
	cp hikup.service $(PACKAGE_NAME)/etc/systemd/system/
	cp control $(PACKAGE_NAME)/DEBIAN/
	dpkg-deb --build $(PACKAGE_NAME)

clean:
	rm -f $(BINARY_NAME)
	rm -rf $(PACKAGE_NAME)
	rm -f $(PACKAGE_NAME).deb

.PHONY: all build package clean
