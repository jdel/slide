package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Log exposes a preconfigured zerolog logger
var Log zerolog.Logger = log.With().Timestamp().Logger().Output(zerolog.ConsoleWriter{
	Out:        os.Stderr,
	TimeFormat: time.RFC3339Nano,
})

// SetLevel parses lvl and sets the logger level
func SetLevel(lvl string) {
	if parsedLevel, err := zerolog.ParseLevel(lvl); err == nil {
		zerolog.SetGlobalLevel(parsedLevel)
		if parsedLevel == zerolog.DebugLevel {
			Log = Log.With().Caller().Logger()
		}
	} else {
		log.Err(err).Send()
	}
}
