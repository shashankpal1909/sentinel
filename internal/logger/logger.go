package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type CustomHandler struct {
	level slog.Level
	out   io.Writer
	mu    *sync.Mutex
	attrs []slog.Attr
}

func NewCustomHandler(out io.Writer, level slog.Level) *CustomHandler {
	return &CustomHandler{
		level: level,
		out:   out,
		mu:    &sync.Mutex{},
	}
}

func (h *CustomHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *CustomHandler) Handle(_ context.Context, r slog.Record) error {
	timeStr := r.Time.Format("2006-01-02T15:04:05.000Z07:00")
	levelStr := r.Level.String()

	var sb strings.Builder
	fmt.Fprintf(&sb, "[%s] [%s] %s", timeStr, levelStr, r.Message)

	for _, a := range h.attrs {
		fmt.Fprintf(&sb, " %s=%v", a.Key, a.Value.Any())
	}

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(&sb, " %s=%v", a.Key, a.Value.Any())
		return true
	})
	sb.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write([]byte(sb.String()))
	return err
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)
	return &CustomHandler{
		level: h.level,
		out:   h.out,
		mu:    h.mu,
		attrs: newAttrs,
	}
}

func (h *CustomHandler) WithGroup(_ string) slog.Handler {
	return h
}

// Init initializes and configures the global custom formatted logger.
func Init(levelStr string) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := NewCustomHandler(os.Stdout, level)
	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}
