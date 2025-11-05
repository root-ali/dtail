package main

import (
	"fmt"
	"io"

	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
	Purple  = "\033[35m"
	Salmon  = "\033[38;5;210m"
)

var Colors = []string{Green, Blue, Yellow, Red, Magenta, Cyan, Gray, White, Purple, Salmon}

type Config struct {
	Output        io.Writer
	ContainerName string
	Color         string
}

func NewLogWriter(cfg Config, stream string) *LogWriter {

	output := zerolog.ConsoleWriter{Out: cfg.Output, NoColor: false}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s ", i)
	}
	output.FormatTimestamp = func(i interface{}) string {
		return fmt.Sprintf("%s |", i)
	}
	output.TimeFormat = time.RFC3339Nano
	output.PartsOrder = []string{"container", "message"}
	output.FormatPartValueByName = func(i interface{}, s string) string {
		var ret string
		switch s {
		case "container":
			log := cfg.Color + i.(string) + " =>" + Reset
			ret = fmt.Sprintf("%s", log)
		}
		return ret
	}
	output.FieldsExclude = []string{"container"}

	log := zerolog.New(output).With().Timestamp().Logger()

	return &LogWriter{
		Logger: &log,
		Config: cfg,
		stream: stream,
	}
}

type LogWriter struct {
	Logger *zerolog.Logger
	Config Config
	stream string
}

func (w LogWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg == "" {
		return len(p), nil
	}

	switch w.stream {
	case "stdout":
		w.Logger.Info().
			Str("container", w.Config.ContainerName).
			Msg(msg)
	case "stderr":
		w.Logger.Error().
			Str("container", w.Config.ContainerName).
			Msg(msg)
	default:
		w.Logger.Warn().
			Str("container", w.Config.ContainerName).
			Msg(fmt.Sprintf("[unknown stream] %s", msg))
	}

	return len(p), nil
}
