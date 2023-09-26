// Package cli implements a colored text handler suitable for command-line interfaces.
package cli

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/fatih/color"
	"github.com/justgook/goplatformer/pkg/gameLogger"
)

// Colors mapping.
var Colors = map[slog.Level]*color.Color{
	slog.LevelDebug: color.New(color.FgWhite),
	slog.LevelInfo:  color.New(color.FgBlue),
	slog.LevelWarn:  color.New(color.FgYellow),
	slog.LevelError: color.New(color.FgRed),
}

// Strings mapping.
var Strings = map[slog.Level]string{
	slog.LevelDebug: "•",
	slog.LevelInfo:  "•",
	slog.LevelWarn:  "•",
	slog.LevelError: "⨯",
}

type Handler struct {
	mu sync.Mutex
	w  io.Writer

	opts *Options

	attrsPrefix []slog.Attr

	groupPrefix string
}

type Options struct {
	DisableColor bool
	slog.HandlerOptions
}

func New(w io.Writer, opts *Options) *Handler {
	var bold = color.New(color.Bold)
	bold.EnableColor()
	return &Handler{w: w, opts: opts}
}

func (h *Handler) Enabled(ctx context.Context, l slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return l >= minLevel
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	var bold = color.New(color.Bold)
	// level time message attributes// get a buffer from the sync pool
	// get a buffer from the sync pool
	buf := gameLogger.NewBuffer()
	defer buf.Free()

	theColor := Colors[r.Level]

	if h.opts.DisableColor {
		theColor.DisableColor()
	} else {
		theColor.EnableColor()
	}

	levelEmoji := Strings[r.Level]
	padding := 4
	coloredLevel := theColor.Sprintf("%s", bold.Sprintf("%*s", padding, levelEmoji))
	buf.WriteString(coloredLevel)
	buf.WriteString(" ")
	buf.WriteString(fmt.Sprintf("%-25s", r.Message))
	buf.WriteString("\t\t")

	// write handler attributes
	if len(h.attrsPrefix) > 0 {
		for _, attr := range h.attrsPrefix {
			h.appendAttr(buf, attr, theColor, h.groupPrefix)
		}
	}

	// write attributes
	if r.NumAttrs() > 0 {
		r.Attrs(func(attr slog.Attr) bool {
			h.appendAttr(buf, attr, theColor, h.groupPrefix)
			return true
		})
	}

	buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) appendAttr(buf *gameLogger.Buffer, attr slog.Attr, theColor *color.Color, groupsPrefix string) {
	buf.Write([]byte(" "))
	if groupsPrefix != "" {
		buf.WriteString(theColor.Sprint(groupsPrefix))
	}
	buf.WriteString(theColor.Sprint(attr.Key))
	buf.Write([]byte("="))

	// needQuote := attr.Value.Kind() != slog.KindInt64
	// if needQuote {
	// 	buf.Write([]byte(`"`))
	// }
	if attr.Value.Kind() != slog.KindGroup {
		buf.Write([]byte(attr.Value.String()))
	} else {
		buf.Write([]byte("{"))
		for _, attr := range attr.Value.Group() {
			h.appendAttr(buf, attr, theColor, groupsPrefix)
		}
		buf.Write([]byte(" }"))
	}

	// if needQuote {
	// 	buf.Write([]byte(`"`))
	// }
}

func (h *Handler) clone() *Handler {
	attrsPrefix := make([]slog.Attr, len(h.attrsPrefix))
	copy(attrsPrefix, h.attrsPrefix)
	return &Handler{w: h.w, opts: h.opts, attrsPrefix: attrsPrefix, groupPrefix: h.groupPrefix}
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	cloned := h.clone()
	cloned.attrsPrefix = append(cloned.attrsPrefix, attrs...)
	return cloned
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	cloned := h.clone()
	cloned.groupPrefix += name + "."
	return cloned
}

// var _ slog.Handler = (*Handler)(nil)
