package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func levelColor(level zerolog.Level, noColor bool) int {
	if noColor {
		return 0
	}
	return zerolog.LevelColors[level]
}

func SetupLogger(verbose, noColor bool) {
	log.Logger = log.Output(zerolog.NewConsoleWriter(
		func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = time.TimeOnly
			w.FormatLevel = func(i interface{}) string {
				fl := ""
				if ll, ok := i.(string); ok {
					level, _ := zerolog.ParseLevel(ll)
					fl = fmt.Sprintf("\x1b[%dm%-5v\x1b[0m", levelColor(level, noColor), strings.ToUpper(ll))
				} else {
					fl = "???"
				}
				return fmt.Sprintf("| %-5s |", fl)
			}
			w.FormatMessage = func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			}
			w.NoColor = noColor
		},
	))

	if verbose {
		log.Logger = log.Level(zerolog.DebugLevel)
	} else {
		log.Logger = log.Level(zerolog.InfoLevel)
	}
}
