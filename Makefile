SOURCES = $(shell find . -name "*.go")

.PHONY: fmt test print-%

all: bump-bedrock

bump-bedrock: $(SOURCES)
	go install

fmt:
	go fmt *.go

test:
	go test -cover

print-%:
	@echo $*=$($*)
