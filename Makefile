SOURCES = $(shell find . ./bedrock -name "*.go")

.PHONY: fmt test print-% coverage

all: bump-bedrock

bump-bedrock: $(SOURCES)
	godep go install

deps:
	godep save -r ./...

fmt:
	go fmt *.go

test:
	go test -cover ./ ./bedrock

coverage: 
	go test -covermode=count -coverprofile=main.out ./
	go test -covermode=count -coverprofile=bedrock.out ./bedrock
	cat main.out bedrock.out | grep -v "mode: count" > all.out
	sed -i '0,/^/s//mode: count/' all.out
	go tool cover -html=all.out

clean:
	rm -f *.out
	rm -rf /tmp/bump-bedrock-test

print-%:
	@echo $*=$($*)
