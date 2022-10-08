BIN := mirakurun_exporter

.PHONY: build
build: $(BIN)

GO_FILES := $(shell find . -type f -name '*.go' -print)

ifdef RELEASE
	GO_LDFLAGS += -w -s
	GO_FLAGS += -a -installsuffix netgo
	GO_BUILD_TAGS := netgo
endif

VERSION := $(shell cat VERSION)
COMMIT_SHA := $(shell git rev-parse --short HEAD)
GO_LDFLAGS += -X 'main.BuildVersion=$(VERSION)' -X 'main.BuildCommitSha=$(COMMIT_SHA)'

$(BIN): $(GO_FILES)
	go build -o $@ -tags=$(GO_BUILD_TAGS) $(GO_FLAGS) -ldflags "$(GO_LDFLAGS)"

.PHONY: clean
clean:
	$(RM) $(BIN)
