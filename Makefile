.PHONY: test

test:
	tinygo test -target wasi -v ./...
