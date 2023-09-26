package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/justgook/goplatformer/pkg/gameLogger/cli"
	"image"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jackmordaunt/icns/v2"
)

func main() {
	handler := cli.New(os.Stderr, &cli.Options{
		HandlerOptions: slog.HandlerOptions{Level: slog.LevelDebug},
	})
	slog.SetDefault(slog.New(handler))
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(2)
	}
}

func run() error {
	var (
		name        = flag.String("name", "My Application", "app name")
		author      = flag.String("author", "Romāns Potašovs", "author")
		version     = flag.String("version", "1.0", "app version")
		identifier  = flag.String("id", "", "bundle identifier")
		icon        = flag.String("icon", "", "icon image file (.icns|.png|.jpg|.jpeg)")
		outPathFlag = flag.String("o", *name+".app", "output location")
	)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return errors.New("missing executable argument")
	}

	binInput := args[0]
	binFile := filepath.Base(binInput)
	outPath := *outPathFlag

	contentsPath := filepath.Join(outPath, "Contents")
	macOSPath := filepath.Join(contentsPath, "MacOS")
	resourcesPath := filepath.Join(contentsPath, "Resources")
	binFilePath := filepath.Join(macOSPath, binFile)

	if err := os.MkdirAll(macOSPath, 0777); err != nil {
		return fmt.Errorf("os.MkdirAll appPath: %w", err)
	}
	fdst, err := os.Create(binFilePath)
	if err != nil {
		return fmt.Errorf("create bin: %w", err)
	}
	defer fdst.Close()
	fsrc, err := os.Open(binInput)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New(binInput + " not found")
		}
		return fmt.Errorf("os.Open: %w", err)
	}
	defer fsrc.Close()
	if _, err := io.Copy(fdst, fsrc); err != nil {
		return fmt.Errorf("copy bin: %w", err)
	}
	if err := exec.Command("chmod", "+x", macOSPath).Run(); err != nil {
		return fmt.Errorf("chmod: %w", err)

	}
	if err := exec.Command("chmod", "+x", binFilePath).Run(); err != nil {
		return fmt.Errorf("chmod %s: %w", binFilePath, err)
	}
	id := *identifier
	if id == "" {
		id = *author + "." + *name
	}
	info := infoListData{
		Name:               *name,
		Executable:         filepath.Join("MacOS", binFile),
		Identifier:         id,
		Version:            *version,
		InfoString:         *name + " by " + *author,
		ShortVersionString: *version,
	}
	if *icon != "" {
		iconPath, err := prepareIcons(*icon, resourcesPath)
		if err != nil {
			return fmt.Errorf("icon: %w", err)
		}
		info.IconFile = filepath.Base(iconPath)
	}
	tpl, err := template.New("template").Parse(infoPlistTemplate)
	if err != nil {
		return fmt.Errorf("infoPlistTemplate: %w", err)
	}
	fplist, err := os.Create(filepath.Join(contentsPath, "Info.plist"))
	if err != nil {
		return fmt.Errorf("create Info.plist: %w", err)
	}
	defer fplist.Close()
	if err := tpl.Execute(fplist, info); err != nil {
		return fmt.Errorf("execute Info.plist template: %w", err)
	}
	if err := os.WriteFile(filepath.Join(contentsPath, "README"), []byte(readme), 0666); err != nil {
		return fmt.Errorf("ioutil.WriteFile: %w", err)
	}
	return nil
}

func prepareIcons(iconPath, resourcesPath string) (string, error) {
	ext := filepath.Ext(strings.ToLower(iconPath))
	fsrc, err := os.Open(iconPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("icon file not found")
		}
		return "", fmt.Errorf("open icon file: %w", err)
	}
	defer fsrc.Close()
	if err2 := os.MkdirAll(resourcesPath, 0777); err2 != nil {
		return "", fmt.Errorf("os.MkdirAll resourcesPath: %w", err2)
	}
	destFile := filepath.Join(resourcesPath, "icon.icns")
	fdst, err := os.Create(destFile)
	if err != nil {
		return "", fmt.Errorf("create icon.icns file: %w", err)
	}
	defer fdst.Close()
	switch ext {
	case ".icns": // just copy the .icns file
		if _, err2 := io.Copy(fdst, fsrc); err2 != nil {
			return destFile, fmt.Errorf("copying %s: %w", iconPath, err2)
		}
	case ".png", ".jpg", ".jpeg", ".gif": // process any images
		srcImg, _, err2 := image.Decode(fsrc)
		if err2 != nil {
			return destFile, fmt.Errorf("decode image: %w", err2)
		}
		if err3 := icns.Encode(fdst, srcImg); err3 != nil {
			return destFile, fmt.Errorf("generate icns file: %w", err3)
		}
	default:
		return destFile, errors.New(ext + " icons not supported")
	}

	return destFile, nil
}

type infoListData struct {
	Name               string
	Executable         string
	Identifier         string
	Version            string
	InfoString         string
	ShortVersionString string
	IconFile           string
}

const infoPlistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>CFBundlePackageType</key>
		<string>APPL</string>
		<key>CFBundleInfoDictionaryVersion</key>
		<string>6.0</string>
		<key>CFBundleName</key>
		<string>{{ .Name }}</string>
		<key>CFBundleExecutable</key>
		<string>{{ .Executable }}</string>
		<key>CFBundleIdentifier</key>
		<string>{{ .Identifier }}</string>
		<key>CFBundleVersion</key>
		<string>{{ .Version }}</string>
		<key>CFBundleGetInfoString</key>
		<string>{{ .InfoString }}</string>
		<key>CFBundleShortVersionString</key>
		<string>{{ .ShortVersionString }}</string>
		{{ if .IconFile -}}
		<key>CFBundleIconFile</key>
		<string>{{ .IconFile }}</string>
		{{- end }}
	</dict>
</plist>
`

const readme = `Made by Romāns Potašovs

`
