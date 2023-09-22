SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifdef V
Q=
WGET:=wget
else
Q=@
MAKEFLAGS += --no-print-directory
WGET:=wget -q --show-progress
endif

define QUIET
	$(if $(V), , $(1))
endef

MKDIR_P ?= mkdir -p
CP ?= cp -f

.DEFAULT_GOAL := all

BUILD_DIR ?= build.nosync
WEB_DIR ?= $(BUILD_DIR)/web
BIN_WIN64_DIR ?= $(BUILD_DIR)/win64
BIN_MAC64_DIR ?= $(BUILD_DIR)/mac64
BIN_MACAARC64_DIR ?= $(BUILD_DIR)/macM1

APP_DIR := ./cmd/game

clean:
	$(Q)git clean -xdf
	$(Q)rm -rf $(BUILD_DIR)
.PHONY: clean

all: web release-win64 release-mac-intel release-mac-arm
.PHONY: all

web: $(WEB_DIR)/game.wasm $(WEB_DIR)/index.html $(WEB_DIR)/wasm_exec.js
.PHONY: web

release-win64: $(BIN_WIN64_DIR)/game.exe
	$(Q)echo use tools to add icons / versions
	$(Q)echo https://github.com/tc-hib/go-winres
	$(Q)echo https://github.com/josephspurrier/goversioninfo

.PHONY: release-win64

release-mac-intel: $(BIN_MAC64_DIR)/game
.PHONY: release-mac-intel

release-mac-arm: $(BIN_MACAARC64_DIR)/game
.PHONY: release-mac-arm

run:
	$(Q)go run $(APP_DIR)
.PHONY: run

update:
	$(Q)go get -u ./...
	$(Q)go mod tidy
.PHONY: update

$(BIN_WIN64_DIR)/game.exe: export GOOS=windows
$(BIN_WIN64_DIR)/game.exe: export GOARCH=amd64
$(BIN_WIN64_DIR)/game.exe:
	$(Q)go build -ldflags="-s -w" -o $@ $(APP_DIR)


$(BIN_MAC64_DIR)/game: export GOOS=darwin
$(BIN_MAC64_DIR)/game: export GOARCH=amd64
$(BIN_MAC64_DIR)/game:
	$(Q)go build -ldflags="-s -w" -o $@ $(APP_DIR)

$(BIN_MACAARC64_DIR)/game: export GOOS=darwin
$(BIN_MACAARC64_DIR)/game: export GOARCH=arm64
$(BIN_MACAARC64_DIR)/game:
	$(Q)go build -ldflags="-s -w" -o $@ $(APP_DIR)

$(WEB_DIR)/game.wasm: export GOOS=js
$(WEB_DIR)/game.wasm: export GOARCH=wasm
$(WEB_DIR)/game.wasm:
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)go build -ldflags="-s -w" -o $@ $(APP_DIR)

$(WEB_DIR)/wasm_exec.js:
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js $@

define INDEX_HTML_CONTENT
<!DOCTYPE html>
<script src="wasm_exec.js"></script>
<script>
// Polyfill
if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

const go = new Go();
WebAssembly.instantiateStreaming(fetch("game.wasm"), go.importObject).then(result => {
    go.run(result.instance);
});
</script>
endef

$(WEB_DIR)/index.html: export INDEX_HTML_CONTENT:=$(INDEX_HTML_CONTENT)
$(WEB_DIR)/index.html:
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)echo "$${INDEX_HTML_CONTENT}" > $@
