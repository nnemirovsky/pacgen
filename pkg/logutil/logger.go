package logutil

import (
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

const (
	pathSeparator = string(filepath.Separator)
	maxPatsCount  = 3
)

var (
	stderr        io.Writer
	Logger        zerolog.Logger
	DiscardLogger = zerolog.New(io.Discard)
	LayerKey      = "layer"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = time.Now().UTC
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	stderr = zerolog.SyncWriter(os.Stderr)
	Logger = NewLogger()
}

func NewLogger() zerolog.Logger {
	var writer io.Writer

	if isatty.IsTerminal(os.Stderr.Fd()) {
		writer = zerolog.ConsoleWriter{
			Out:          stderr,
			TimeFormat:   time.RFC3339Nano,
			FormatCaller: formatCaller,
		}
	} else {
		writer = stderr
	}

	logger := zerolog.New(writer).Level(zerolog.InfoLevel).With().
		Timestamp().
		Caller().
		Logger()

	return logger
}

func WithLayer[T any](logger zerolog.Logger) zerolog.Logger {
	return logger.With().Str(LayerKey, getTypeNameWithPkg[T]()).Logger()
}

func getTypeNameWithPkg[T any]() string {
	var v T
	t := reflect.TypeOf(v)
	pkgPath := strings.Split(t.PkgPath(), "/")
	return pkgPath[len(pkgPath)-1] + "." + t.Name()
}

func formatCaller(c any) string {
	caller := strings.TrimPrefix(c.(string), pathSeparator)
	parts := strings.Split(caller, pathSeparator)
	if len(parts) > maxPatsCount {
		parts = parts[len(parts)-maxPatsCount:]
	}
	return strings.Join(parts, pathSeparator)
}
