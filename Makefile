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
ASSET_DIR ?= asset
WEB_DIR ?= $(BUILD_DIR)/web
BIN_WIN64_DIR ?= $(BUILD_DIR)/win64
BIN_MAC64_DIR ?= $(BUILD_DIR)/mac64
BIN_MACAARC64_DIR ?= $(BUILD_DIR)/macM1

SRC_LEVEL = $(wildcard $(ASSET_DIR)/*.ldtk)
SRC_LEVEL += $(wildcard $(ASSET_DIR)/*/*.ldtk)
SRC_LEVEL += $(wildcard $(ASSET_DIR)/*/*/*.ldtk)

SRC_SPRITE = $(wildcard $(ASSET_DIR)/*.aseprite)
SRC_SPRITE += $(wildcard $(ASSET_DIR)/*/*.aseprite)
SRC_SPRITE += $(wildcard $(ASSET_DIR)/*/*/*.aseprite)

RESOURCES := $(SRC_LEVEL:.ldtk=.level)
RESOURCES += $(SRC_SPRITE:.aseprite=.sprite)
RESOURCES := $(subst $(ASSET_DIR),$(BUILD_DIR),$(RESOURCES))

APP_DIR := ./cmd/game


LDFLAGS := -s -w

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

release-mac-intel: $(BUILD_DIR)/MyGameX86.app
.PHONY: release-mac-intel

release-mac-arm: $(BUILD_DIR)/MyGameM1.app
.PHONY: release-mac-arm

APPIFY := go run github.com/justgook/goplatformer/cmd/appify

$(BUILD_DIR)/MyGameM1.app: $(BIN_MACAARC64_DIR)/game
	$(Q)$(APPIFY) -name "My Super Game" -icon ./asset/bundle/icon.png -o $@ $<

$(BUILD_DIR)/MyGameX86.app: $(BIN_MAC64_DIR)/game
	$(Q)$(APPIFY) -name "My Super Game" -icon ./asset/bundle/icon.png -o $@ $<


run: #$(RESOURCES)
	$(Q)go run $(APP_DIR)
.PHONY: run

update:
	$(Q)go get -u ./...
	$(Q)go mod tidy
.PHONY: update


$(BIN_WIN64_DIR)/game.exe: export GOOS=windows
$(BIN_WIN64_DIR)/game.exe: export GOARCH=amd64
$(BIN_WIN64_DIR)/game.exe: $(RESOURCES)
	$(Q)go build -ldflags="$(LDFLAGS)" -o $@ $(APP_DIR)


$(BIN_MAC64_DIR)/game: export GOOS=darwin
$(BIN_MAC64_DIR)/game: export GOARCH=amd64
$(BIN_MAC64_DIR)/game: $(RESOURCES)
	$(Q)go build -ldflags="$(LDFLAGS)" -o $@ $(APP_DIR)

$(BIN_MACAARC64_DIR)/game: export GOOS=darwin
$(BIN_MACAARC64_DIR)/game: export GOARCH=arm64
$(BIN_MACAARC64_DIR)/game: $(RESOURCES)
	$(Q)go build -ldflags="$(LDFLAGS)" -o $@ $(APP_DIR)

$(WEB_DIR)/game.wasm: export GOOS=js
$(WEB_DIR)/game.wasm: export GOARCH=wasm
$(WEB_DIR)/game.wasm: #$(RESOURCES)
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)go build -ldflags="$(LDFLAGS)" -o $@ $(APP_DIR)

$(WEB_DIR)/wasm_exec.js:
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js $@

define INDEX_HTML_CONTENT
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Permissions-Policy" content="interest-cohort=()">
</head>
<body>
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
document.body.querySelectorAll("script").forEach(a => a.remove())
</script>
</body>
</html>
endef

$(WEB_DIR)/index.html: export INDEX_HTML_CONTENT:=$(INDEX_HTML_CONTENT)
$(WEB_DIR)/index.html:
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)echo "$${INDEX_HTML_CONTENT}" > $@


ASEPRITE = /Applications/Aseprite.app/Contents/MacOS/aseprite
# ifneq ("$(wildcard $(ASEPRITE))","")
#
# else
# 	# ASEPRITE="C:\Program Files\Aseprite\aseprite.exe"
# 	$(error please install Aseprite)
# endif

resources: $(RESOURCES)

$(BUILD_DIR)/%.level: $(ASSET_DIR)/%.ldtk
	$(Q)echo building $@ from $<
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)go run ./cmd/level -o $@ $<

$(BUILD_DIR)/%.png $(BUILD_DIR)/%.json &: $(ASSET_DIR)/%.aseprite
	$(Q)echo building $@ from $<
	$(Q)$(MKDIR_P) $(dir $@)
	$(ASEPRITE) -b \
		--data $(BUILD_DIR)/$*.json \
		--format json-hash \
		--sheet $(BUILD_DIR)/$*.png \
		--sheet-type packed \
		--sheet-pack \
		--split-layers \
		--ignore-empty \
		--merge-duplicates \
		--trim \
		--trim-sprite \
		--extrude \
		--filename-format '{layer} {frame}' \
		--tagname-format {tag} \
		--list-layers \
		--list-tags \
		$<

$(BUILD_DIR)/%.sprite: $(BUILD_DIR)/%.png $(BUILD_DIR)/%.json
	$(Q)echo building $@ from $<
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)go run ./cmd/sprite \
		-o $@ \
		--data $(BUILD_DIR)/$*.json \
		--sprite $(BUILD_DIR)/$*.png

