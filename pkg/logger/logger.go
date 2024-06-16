package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type logger struct {
	*zerolog.Logger
}

type Logger interface {
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Fatal() *zerolog.Event
}

func New() (Logger, error) {
	l := zerolog.New(os.Stdout)

	return &logger{
		&l,
	}, nil
}
