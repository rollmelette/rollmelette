.PHONY: all
all: | lint build test

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -p 1 ./...

.PHONY: cov
cov:
	go test -p 1 -race -coverpkg=. -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func coverage.txt

.PHONY: gen
gen:
	go generate ./...

.PHONY: check-gen
check-gen: gen
	git diff --quiet

.PHONY: run
run:
	go run ./examples/ $(EXAMPLE)

.PHONY: clean
clean:
	rm -f coverage.txt
