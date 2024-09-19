package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/padok-team/git-secret-scanner/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cmd.Execute()
}

func init() {
	log.Logger = log.Output(zerolog.NewConsoleWriter(
		func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = time.TimeOnly
			w.FormatLevel = func(i interface{}) string {
				ll := fmt.Sprintf("%s", i)
				if level, err := zerolog.ParseLevel(ll); err == nil {
					ll = fmt.Sprintf("\x1b[%dm%-5v\x1b[0m", zerolog.LevelColors[level], strings.ToUpper(ll))
				}
				return fmt.Sprintf("| %-5s |", ll)
			}
			w.FormatMessage = func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			}
		},
	))
	log.Logger = log.Level(zerolog.InfoLevel)
}
