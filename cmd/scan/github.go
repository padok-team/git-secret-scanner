package scan

import (
	"context"
	"os"

	"github.com/padok-team/git-secret-scanner/internal/scan"
	"github.com/padok-team/git-secret-scanner/internal/scm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const TokenEnvVarGithub string = "GITHUB_TOKEN"

var GithubCmd = &cobra.Command{
	Use:     "github",
	GroupID: "scan",
	Short:   "Scan for secrets in a GitHub organization",
	Args:    cobra.NoArgs,
	PreRun:  preRun,
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv(TokenEnvVarGithub)
		if token == "" {
			log.Fatal().Msgf("missing environment variable %q", TokenEnvVarGithub)
		}

		github, err := scm.NewGithubScm(&scmConfig, token)
		if err != nil {
			log.Fatal().Msgf("failed to initialize github client: %s", err)
		}

		log.Debug().
			Str("org", scmConfig.Org).
			Str("server", scmConfig.Server).
			Str("git_protocol", scmConfig.GitProtocol.String()).
			Str("visibility", scmConfig.Visiblity.String()).
			Bool("include_archived", scmConfig.IncludeArchived).
			Msg("parsed github config")

		err = scan.Scan(context.Background(), github, scanArgs)
		if err != nil {
			log.Fatal().Msgf("an error occured during the scan: %s", err)
		}
	},
}

func init() {
	GithubCmd.Flags().StringVarP(&scmConfig.Org, "org", "o", "", "Organization to scan (required)")
	GithubCmd.MarkFlagRequired("org") //nolint:errcheck

	registerCommonFlags(GithubCmd)
}
