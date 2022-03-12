GO=go

.PHONY: test

test:
	$(GO) test ./...

bench:
	$(GO) test -bench=.
