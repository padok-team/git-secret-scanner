package scm

import (
	"fmt"

	"github.com/padok-team/git-secret-scanner/internal/scan"
	"github.com/padok-team/git-secret-scanner/internal/scm"
	"github.com/padok-team/git-secret-scanner/internal/scm/git"
	"github.com/padok-team/git-secret-scanner/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	ScannerBinaryTrufflehog string = "trufflehog"
	ScannerBinaryGitleaks   string = "gitleaks"
)

var (
	scanArgs  scan.ScanArgs
	scmConfig scm.ScmConfig

	noArchived bool
	sshClone   bool
	filesOnly  bool
	verbose    bool
)

func preRun(cmd *cobra.Command, args []string) {
	if !utils.CommandExists(ScannerBinaryTrufflehog) {
		log.Fatal().Msgf("executable %q not found in PATH", ScannerBinaryTrufflehog)
	}
	if !utils.CommandExists(ScannerBinaryGitleaks) {
		log.Fatal().Msgf("executable %q not found in PATH", ScannerBinaryGitleaks)
	}

	if verbose {
		log.Logger = log.Level(zerolog.DebugLevel)
	}

	// parse flags
	scmConfig.IncludeArchived = !noArchived
	if sshClone {
		scmConfig.GitProtocol = git.ProtocolSSH
	}
	if filesOnly {
		scanArgs.ScanType = scan.ScanTypeFilesOnly
	}

	log.Debug().
		Str("scan_type", scanArgs.ScanType.String()).
		Str("report_path", scanArgs.ReportPath).
		Str("fingerprints_ignore_path", scanArgs.FingerprintsIgnorePath).
		Str("baseline_path", scanArgs.BaselinePath).
		Int("max_concurrency", scanArgs.MaxConcurrency).
		Msg("parsed scan args")
}

func registerCommonFlags(cmd *cobra.Command) {
	// scm config (org is not provided here as the flag name changes between scm)
	cmd.Flags().VarP(&scmConfig.Visiblity, "visibility", "v", "Visibility of repositories to scan")
	cmd.Flags().StringVar(&scmConfig.Server, "server", "", "Hostname of the server")
	cmd.Flags().BoolVar(&noArchived, "no-archived", false, "Skip archived repositories")
	cmd.Flags().BoolVar(&sshClone, "ssh-clone", false, "Use SSH to clone repositories instead of HTTPS")

	// scan args
	cmd.Flags().StringVarP(&scanArgs.ReportPath, "report-path", "r", "report.csv", "Path to the CSV report file to generate")
	cmd.Flags().StringVarP(&scanArgs.FingerprintsIgnorePath, "fingerprints-ignore-path", "i", "", "Path to file with newline separated fingerprints (SHA-256) of secrets to ignore during the scan")
	cmd.Flags().StringVarP(&scanArgs.BaselinePath, "baseline-path", "b", "", "Path to the CSV report to use as baseline for the scan")
	cmd.Flags().IntVar(&scanArgs.MaxConcurrency, "max-concurrency", 5, "Maximum number of concurrent workers")
	cmd.Flags().BoolVarP(&filesOnly, "files-only", "f", false, "Only run the scan on the files of the default branch")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Show verbose output")

	// help flag
	cmd.Flags().BoolP("help", "h", false, fmt.Sprintf("Help for command %s", cmd.Name()))
}

func AddScmCommands(cmd *cobra.Command) {
	cmd.AddCommand(githubCmd)
	cmd.AddCommand(gitlabCmd)
}
