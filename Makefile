SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=
OS=$(shell uname -s)

build:
	go generate ./...
	go build -o pd-report
.PHONY: build

test:
	go test $(TEST_OPTIONS) -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.out $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m
.PHONY: test

# Run all the tests and opens the coverage report
cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

# gofmt and goimports all go files
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: fmt

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi
.PHONY: vet

embed-assets:
	rice -v -i ./configuration embed-go
.PHONY: embed-assets

ci: test vet
.PHONY: ci

# Show to-do items per file.
todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--exclude-dir=.idea \
		--exclude-dir=bin \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E 'TODO' .
.PHONY: todo


.DEFAULT_GOAL := build
