package scan

import (
	"context"
	"os"

	"github.com/padok-team/git-secret-scanner/internal/scan"
	"github.com/padok-team/git-secret-scanner/internal/scm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const TokenEnvVarGitlab string = "GITLAB_TOKEN"

var GitlabCmd = &cobra.Command{
	Use:     "gitlab",
	GroupID: "scan",
	Short:   "Scan for secrets in a Gitlab group",
	Args:    cobra.NoArgs,
	PreRun:  preRun,
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv(TokenEnvVarGitlab)
		if token == "" {
			log.Fatal().Msgf("missing environment variable %q", TokenEnvVarGitlab)
		}

		gitlab, err := scm.NewGitlabScm(&scmConfig, token)
		if err != nil {
			log.Fatal().Msgf("failed to initialize github client: %s", err)
		}

		log.Debug().
			Str("group", scmConfig.Org).
			Str("server", scmConfig.Server).
			Str("git_protocol", scmConfig.GitProtocol.String()).
			Str("visibility", scmConfig.Visiblity.String()).
			Bool("include_archived", scmConfig.IncludeArchived).
			Msg("parsed gitlab config")

		err = scan.Scan(context.Background(), gitlab, scanArgs)
		if err != nil {
			log.Fatal().Msgf("an error occured during the scan: %s", err)
		}
	},
}

func init() {
	GitlabCmd.Flags().StringVarP(&scmConfig.Org, "group", "g", "", "Group to scan (required)")
	GitlabCmd.MarkFlagRequired("group") //nolint:errcheck

	registerCommonFlags(GitlabCmd)
}
