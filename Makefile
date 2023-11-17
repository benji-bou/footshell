.PHONY: all build display-plugins modifier-plugins clean

all: build display-plugins modifier-plugins

build:
	go build -o footshell ./cmd/main.go

display-plugins:
	@echo "Building display plugins..."
	@for dir in $(wildcard ./plugins/display/*); do \
		go build -buildmode=plugin -o $$dir/$$(basename $$dir).so $$dir/$$(basename $$dir).go; \
	done

modifier-plugins:
	@echo "Building modifier plugins..."
	@for dir in $(wildcard ./plugins/modifiers/*); do \
		go build -buildmode=plugin -o $$dir/$$(basename $$dir).so $$dir/$$(basename $$dir).go; \
	done

clean:
	@echo "Cleaning up..."
	rm -f footshell
	find ./plugins/display -name '*.so' -delete
	find ./plugins/modifiers -name '*.so' -delete