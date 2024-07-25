package logger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"

	m "about-me-bot/internal/models"
	"github.com/google/uuid"
)

const LevelFatal slog.Level = slog.Level(12)

type CustomHandler struct {
	Default slog.Handler
}

func (h *CustomHandler) ToLogger() *slog.Logger {
	return slog.New(h)
}

func (h *CustomHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return l >= slog.LevelDebug
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	var (
		prefix    string
		location  string
		timestamp string = "[" + r.Time.Format("2006-01-02 15:04:05.000") + "]"
		trace            = fmt.Sprintf("[trace-id:%s]", uuid.New().String())
	)
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		SimpleLogger.Error("could not get runtime info on the source")
		file = "unknown"
		line = 0
	}
	function := runtime.FuncForPC(pc).Name()

	splitFunc := strings.Split(function, "/")
	function = splitFunc[len(splitFunc)-1]

	if file != "unknown" {
		i := strings.Index(file, "about-me-bot")
		location = "./" + file[i:] + "/" + function + ":" + fmt.Sprint(line)
	} else {
		location = "./" + file + "/" + function + ":" + fmt.Sprint(line)
	}

	switch r.Level {
	case slog.LevelDebug:
		prefix = m.DbgColor("DEBUG")
	case slog.LevelInfo:
		prefix = m.Cyan("INFO")
	case slog.LevelWarn:
		prefix = m.Yellow("WARNING")
	case slog.LevelError:
		prefix = m.Red("ERROR")
	case LevelFatal:
		defer os.Exit(1)
		prefix = m.DarkRed("FATAL")
	default:
		h.Default.Handle(ctx, r)
		return errors.New("unknown level")
	}

	prefix += " "

	log.New(os.Stderr, prefix, 0).Println(trace, location, timestamp, r.Message)

	return nil
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{h.Default.WithAttrs(attrs)}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{h.Default.WithGroup(name)}
}

var (
	DefaultHandler = slog.Default().Handler()

	SimpleLogger = slog.New(&CustomHandler{Default: DefaultHandler})
)
