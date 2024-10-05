package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
)

type Handler struct {
	h                slog.Handler
	r                func([]string, slog.Attr) slog.Attr
	b                *bytes.Buffer
	m                *sync.Mutex
	writer           io.Writer
	colorize         bool
	outputEmptyAttrs bool
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, r: h.r, m: h.m, writer: h.writer, colorize: h.colorize}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, r: h.r, m: h.m, writer: h.writer, colorize: h.colorize}
}

func (h *Handler) computeAttrs(
	ctx context.Context,
	r slog.Record,
) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}
	return attrs, nil
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	colorize := func(code int, value string) string {
		return value
	}
	if h.colorize {
		colorize = colorizer
	}

	var level string
	levelAttr := slog.Attr{
		Key:   slog.LevelKey,
		Value: slog.AnyValue(r.Level),
	}
	if h.r != nil {
		levelAttr = h.r([]string{}, levelAttr)
	}

	if !levelAttr.Equal(slog.Attr{}) {
		level = levelAttr.Value.String() + ":"

		if r.Level <= slog.LevelDebug {
			level = colorize(lightGray, level)
		} else if r.Level <= slog.LevelInfo {
			level = colorize(cyan, level)
		} else if r.Level < slog.LevelWarn {
			level = colorize(lightBlue, level)
		} else if r.Level < slog.LevelError {
			level = colorize(lightYellow, level)
		} else if r.Level <= slog.LevelError+1 {
			level = colorize(lightRed, level)
		} else if r.Level > slog.LevelError+1 {
			level = colorize(lightMagenta, level)
		}
	}

	var timestamp string
	timeAttr := slog.Attr{
		Key:   slog.TimeKey,
		Value: slog.StringValue(r.Time.Format(timeFormat)),
	}
	if h.r != nil {
		timeAttr = h.r([]string{}, timeAttr)
	}
	if !timeAttr.Equal(slog.Attr{}) {
		timestamp = colorize(lightGray, timeAttr.Value.String())
	}

	var msg string
	msgAttr := slog.Attr{
		Key:   slog.MessageKey,
		Value: slog.StringValue(r.Message),
	}
	if h.r != nil {
		msgAttr = h.r([]string{}, msgAttr)
	}
	if !msgAttr.Equal(slog.Attr{}) {
		msg = colorize(white, msgAttr.Value.String())
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	var attrsAsBytes []byte
	if h.outputEmptyAttrs || len(attrs) > 0 {
		attrsAsBytes, err = json.MarshalIndent(attrs, "", "  ")
		if err != nil {
			return fmt.Errorf("error when marshaling attrs: %w", err)
		}
	}

	out := strings.Builder{}
	if len(timestamp) > 0 {
		out.WriteString(timestamp)
		out.WriteString(" ")
	}
	if len(level) > 0 {
		out.WriteString(level)
		out.WriteString(" ")
	}
	if len(msg) > 0 {
		out.WriteString(msg)
		out.WriteString(" ")
	}
	if len(attrsAsBytes) > 0 {
		out.WriteString(colorize(darkGray, string(attrsAsBytes)))
	}

	_, err = io.WriteString(h.writer, out.String()+"\n")
	if err != nil {
		return err
	}

	return nil
}
