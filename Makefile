GOCMD=go
GOTEST=$(GOCMD) test -v
GOBUILD=$(GOCMD) build
BINARY_NAME=ipa-medit

all: build deploy

test:
	$(GOTEST) ./pkg/*

clean:
	rm $(BINARY_NAME)
