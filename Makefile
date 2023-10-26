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

SRC_TILESET = $(wildcard $(ASSET_DIR)/*.tileset.aseprite)
SRC_TILESET += $(wildcard $(ASSET_DIR)/*/*.tileset.aseprite)
SRC_TILESET += $(wildcard $(ASSET_DIR)/*/*/*.tileset.aseprite)

SRC_SPRITE = $(wildcard $(ASSET_DIR)/*.sprite.aseprite)
SRC_SPRITE += $(wildcard $(ASSET_DIR)/*/*.sprite.aseprite)
SRC_SPRITE += $(wildcard $(ASSET_DIR)/*/*/*.sprite.aseprite)

SRC_CHAR = $(wildcard $(ASSET_DIR)/*.char.aseprite)
SRC_CHAR += $(wildcard $(ASSET_DIR)/*/*.char.aseprite)
SRC_CHAR += $(wildcard $(ASSET_DIR)/*/*/*.char.aseprite)

RESOURCES := $(SRC_LEVEL:.ldtk=.level)
RESOURCES += $(SRC_TILESET:.tileset.aseprite=.tileset.png)
RESOURCES += $(SRC_SPRITE:.sprite.aseprite=.sprite)
RESOURCES += $(SRC_CHAR:.char.aseprite=.char)
RESOURCES := $(subst $(ASSET_DIR),$(BUILD_DIR),$(RESOURCES))

APP_DIR := ./cmd/game

SYS_GOOS := $(shell go env GOOS)
SYS_GOARCH := $(shell go env GOARCH)

TAGS ?= development
ifeq ($(BUILD_DIR),build)
  TAGS=production
endif

LDFLAGS := -s -w

clean:
	$(Q)git ls-files -o | xargs trash
	$(Q)rm -rf $(BUILD_DIR)
.PHONY: clean

all: web release-win64 release-mac-intel release-mac-arm
.PHONY: all

develop:
	ls -d pkg/**/* | entr -r -s "make run"
.PHONY: develop

test:
	$(Q)GOOS=js GOARCH=wasm go test ./... -exec="$(shell go env GOROOT)/misc/wasm/go_js_wasm_exec"
.PHONY: test

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



APPIFY := GOOS=$(SYS_GOOS) GOARCH=$(SYS_GOARCH) go run github.com/justgook/goplatformer/cmd/appify

$(BUILD_DIR)/MyGameM1.app: $(BIN_MACAARC64_DIR)/game
	$(Q)$(APPIFY) -name "My Super Game" -icon ./asset/bundle/icon.png -o $@ $<

$(BUILD_DIR)/MyGameX86.app: $(BIN_MAC64_DIR)/game
	$(Q)$(APPIFY) -name "My Super Game" -icon ./asset/bundle/icon.png -o $@ $<


run: $(RESOURCES)
	$(Q)go run -tags $(TAGS) $(APP_DIR)
.PHONY: run

update:
	$(Q)go get -u ./...
	$(Q)go mod tidy
.PHONY: update


$(BIN_WIN64_DIR)/game.exe: export GOOS=windows
$(BIN_WIN64_DIR)/game.exe: export GOARCH=amd64
$(BIN_WIN64_DIR)/game.exe: $(RESOURCES)
	$(Q)echo ... building $@ from $<
	$(Q)go build -tags $(TAGS) -ldflags="$(LDFLAGS)" -o $@ $(APP_DIR)


$(BIN_MAC64_DIR)/game: export GOOS=darwin
$(BIN_MAC64_DIR)/game: export GOARCH=amd64
$(BIN_MACAARC64_DIR)/game: export CGO_CFLAGS=-mmacosx-version-min=10.12
$(BIN_MACAARC64_DIR)/game: export CGO_LDFLAGS=-mmacosx-version-min=10.12
$(BIN_MAC64_DIR)/game: $(RESOURCES)
	$(Q)echo ... building $@ from $<
	$(Q)go build -tags $(TAGS) -ldflags="$(LDFLAGS)" -o $@ $(APP_DIR)

$(BIN_MACAARC64_DIR)/game: export GOOS=darwin
$(BIN_MACAARC64_DIR)/game: export GOARCH=arm64
$(BIN_MACAARC64_DIR)/game: export CGO_CFLAGS=-mmacosx-version-min=10.12
$(BIN_MACAARC64_DIR)/game: export CGO_LDFLAGS=-mmacosx-version-min=10.12
$(BIN_MACAARC64_DIR)/game: $(RESOURCES)
	$(Q)echo ... building $@ from $<
	$(Q)go build -tags $(TAGS) -ldflags="$(LDFLAGS)" -o $@ $(APP_DIR)

$(WEB_DIR)/game.wasm: export GOOS=js
$(WEB_DIR)/game.wasm: export GOARCH=wasm
$(WEB_DIR)/game.wasm: $(RESOURCES)
	$(Q)echo ... building $@ from $<
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)go build -tags $(TAGS) -ldflags="$(LDFLAGS)" -o $@ $(APP_DIR)

$(WEB_DIR)/wasm_exec.js:
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js $@

define INDEX_HTML_CONTENT
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Permissions-Policy" content="interest-cohort=()">
</head>
<body style="height:100vh;display:flex;justify-content:center;align-items:center">
<progress id="progress" value="0" max="1"></progress>
<script src="wasm_exec.js"></script>
<script>
async function fetchWithProgress(path, progress) {
	const response = await fetch(path)
	// May be incorrect if compressed
	const contentLength = response.headers.get("Content-Length")
	const total = parseInt(contentLength, 10)

	let bytesLoaded = 0
	const ts = new TransformStream({
		transform(chunk, ctrl) {
			bytesLoaded += chunk.byteLength
			progress(bytesLoaded / total)
			ctrl.enqueue(chunk)
		},
	})

	return new Response(response.body.pipeThrough(ts), response)
}
// Polyfill
if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}
const go = new Go();
const aaa = document.getElementById("progress")
//WebAssembly.instantiateStreaming(fetchWithProgress("game.wasm", (a) => aaa.value = a), go.importObject).then(result => {
WebAssembly.instantiateStreaming(fetch("game.wasm"), go.importObject).then(result => {
    go.run(result.instance);
    aaa.remove()
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
ifeq ("$(wildcard $(ASEPRITE))","")
	ASEPRITE=aseprite/aseprite
# else
# 	# ASEPRITE="C:\Program Files\Aseprite\aseprite.exe"
# 	$(error please install Aseprite)
endif

resources: $(RESOURCES)

$(BUILD_DIR)/%.level: export GOOS=$(SYS_GOOS)
$(BUILD_DIR)/%.level: export GOARCH=$(SYS_GOARCH)
$(BUILD_DIR)/%.level: $(ASSET_DIR)/%.ldtk $(BUILD_DIR)/%.tileset.png
	$(Q)echo ... building $@ from $<
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)go run ./cmd/level \
		-o $@ \
		--tileset $(BUILD_DIR)/$*.tileset.png \
		--level $(ASSET_DIR)/$*.ldtk

$(BUILD_DIR)/%.sprite: export GOOS=$(SYS_GOOS)
$(BUILD_DIR)/%.sprite: export GOARCH=$(SYS_GOARCH)
$(BUILD_DIR)/%.sprite: $(BUILD_DIR)/%.sprite.png $(BUILD_DIR)/%.sprite.json
	$(Q)echo ... building $@ from $<
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)go run ./cmd/sprite \
		-o $@ \
		--data $(BUILD_DIR)/$*.sprite.json \
		--sprite $(BUILD_DIR)/$*.sprite.png

$(BUILD_DIR)/%.char: export GOOS=$(SYS_GOOS)
$(BUILD_DIR)/%.char: export GOARCH=$(SYS_GOARCH)
$(BUILD_DIR)/%.char: $(BUILD_DIR)/%.char.png $(BUILD_DIR)/%.char.json
	$(Q)echo ... building $@ from $<
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)go run ./cmd/char \
		-o $@ \
		--data $(BUILD_DIR)/$*.char.json \
		--sprite $(BUILD_DIR)/$*.char.png

$(BUILD_DIR)/%.tileset.png: $(ASSET_DIR)/%.tileset.aseprite $(ASEPRITE)
	$(Q)echo ... building $@ from $<
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)$(ASEPRITE) -b \
		$< \
		--save-as $@ \
		--oneframe

$(BUILD_DIR)/%.sprite.png $(BUILD_DIR)/%.sprite.json &: $(ASSET_DIR)/%.sprite.aseprite $(ASEPRITE)
	$(Q)echo ... building $@ from $<
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)$(ASEPRITE) -b \
		--data $(BUILD_DIR)/$*.sprite.json \
		--format json-hash \
		--sheet $(BUILD_DIR)/$*.sprite.png \
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

$(BUILD_DIR)/%.char.png $(BUILD_DIR)/%.char.json &: $(ASSET_DIR)/%.char.aseprite $(ASEPRITE)
	$(Q)echo ... building $@ from $<
	$(Q)$(MKDIR_P) $(dir $@)
	$(Q)$(ASEPRITE) -b \
		--data $(BUILD_DIR)/$*.char.json \
		--format json-hash \
		--filename-format '{tag} {frame}' \
		--sheet $(BUILD_DIR)/$*.char.png \
		--sheet-type packed \
		--ignore-empty \
		--merge-duplicates \
		--extrude \
		--tagname-format {tag} \
		--list-tags \
		$<

aseprite/aseprite:
	$(Q)$(eval TMP := $(shell mktemp -d))
	$(Q)git clone --recursive https://github.com/aseprite/aseprite.git $(TMP)/aseprite
	$(Q)cd $(TMP)/aseprite && cmake -S . -B build -G Ninja \
		-DCMAKE_BUILD_TYPE=Release \
		-DCMAKE_OSX_DEPLOYMENT_TARGET=10.9 \
		-DENABLE_TESTS=off \
		-DENABLE_UI=off \
		-DENABLE_CCACHE=off \
		-DPNG_ARM_NEON=off
	$(Q)cd $(TMP)/aseprite/build && ninja
	$(Q)$(MKDIR_P) "$(CURDIR)/aseprite"
	$(Q)cp -R $(TMP)/aseprite/build/bin/aseprite $(TMP)/aseprite/build/bin/data "$(CURDIR)/aseprite/"
	$(Q)rm -rf $(TMP)
