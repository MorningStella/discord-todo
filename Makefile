EXECUTABLE := discord-bot
SOURCES ?= $(shell find . -name "*.go" -type f)
GO ?= go

build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	mkdir -p bin
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s \
	-w $(LDFLAGS)' -o bin/$@ ./cmd/main.go
    
# install air command
.PHONY: air
air:
	@hash air > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install github.com/cosmtrek/air@latest; \
	fi

# run air
.PHONY: dev
dev: air
	air --build.cmd "make" --build.bin bin/discord-bot

# run the compiled executable
.PHONY: run
run: build
	./bin/$(EXECUTABLE)