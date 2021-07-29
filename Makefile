GOCMD=go
GOTEST=$(GOCMD) test -v
GOBUILD=$(GOCMD) build
BINARY_NAME=ipa-medit

test:
	$(GOTEST) ./pkg/*

build:
	$(GOBUILD) -o $(BINARY_NAME)
	./scripts/codesign.sh

clean:
	rm $(BINARY_NAME)
