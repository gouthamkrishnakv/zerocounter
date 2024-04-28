SHELL = /bin/sh

.PHONY: build clean run watch

# main source file
SOURCE = cmd/main.go

# Go executable (if you require to change)
GO_EXE= go

# Put your output executable name here
BINARY_EXE = zerocounter

build:
	$(GO_EXE) build -o $(BINARY_EXE) $(SOURCE)

run: build
	./$(BINARY_EXE)

watch:
	reflex -s -r '\.go$$' $(GO_EXE) run $(SOURCE)
