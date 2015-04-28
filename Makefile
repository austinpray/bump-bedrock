SOURCES = $(shell find . -name "*.go")

all: bump-bedrock

bump-bedrock: $(SOURCES)
	go install

fmt:
	go fmt *.go

print-%:
	@echo $*=$($*)
