GO?=go

.PHONY: all
all: dictionary/dictionary.txt run

.PHONY: run
run:
	$(GO) run ./cmd/wordworlds/main.go

dictionary/dictionary.txt:
	# NOTE: install aspell if not already.
	aspell -d en dump master | aspell -l en expand | tr " " "\n" | grep -E "^[a-z]{3,}$$" | sort | uniq >$@
