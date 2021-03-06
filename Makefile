GO=go
GOFMT=gofmt
DELETE=rm
BINARY=yabean
BUILD_BINARY=bin/$(BINARY)
# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
# current git version short-hash
VER = $(shell git rev-parse --short HEAD)

info:
	@echo " target         ⾖ Description.                                    "
	@echo " ----------------------------------------------------------------- "
	@echo
	@echo " build          generate a local build ⇨ $(BUILD_BINARY)          "
	@echo " clean          clean up bin/ & go test cache                      "
	@echo " fmt            format go code files using go fmt                  "
	@echo " release/darwin generate a darwin target build                     "
	@echo " release/linux  generate a linux target build                      "
	@echo " tidy           clean up go module file                            "
	@echo
	@echo " ------------------------------------------------------------------"

build: clean fmt
	$(GO) build -o $(BUILD_BINARY) -v main.go


.PHONY: clean
clean:
	$(DELETE) -rf bin/
	$(GO) clean -cache


.PHONY: fmt
fmt:
	$(GOFMT) -l -w $(SRC)


release/%: clean fmt
	@echo "build no race on alpine. https://github.com/golang/go/issues/14481"
	@echo "build GOOS: $(subst release/,,$@) & GOARCH: amd64"
	GOOS=$(subst release/,,$@) GOARCH=amd64 $(GO) build -o bin/$(subst release/,,$@)/$(BINARY) -v main.go

.PHONY: tidy
tidy:
	$(GO) mod tidy