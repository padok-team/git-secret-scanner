package scan

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var ignoredFingerprints map[string][]string

func SetIgnoredFingerprints(fps map[string][]string) {
	ignoredFingerprints = fps
}

func GetRepoIgnoredFingerprints(repository string) []string {
	fps, ok := ignoredFingerprints[repository]
	if !ok {
		return make([]string, 0)
	}
	return fps
}

func LoadIgnoredFingerprints(path string) error {
	if path != "" {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		r, err := regexp.Compile(`^(.+:){3}[0-9]+$`)
		if err != nil {
			return err
		}

		fps := make(map[string][]string, 0)
		num := 0

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			t := scanner.Text()
			if !r.MatchString(t) {
				return errors.New("wrongly formatted fingerprint in file")
			}
			repository := strings.SplitN(t, ":", 4)[0]
			_, ok := fps[repository]
			if !ok {
				fps[repository] = make([]string, 0)
			}
			fps[repository] = append(fps[repository], t)
			num += 1
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		SetIgnoredFingerprints(fps)

		log.Debug().
			Int("num", num).
			Msg("ignored fingerprints loaded")
	}

	return nil
}
