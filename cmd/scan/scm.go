package scan

import (
	"fmt"
	"os/exec"

	"github.com/padok-team/git-secret-scanner/internal/logging"
	"github.com/padok-team/git-secret-scanner/internal/scan"
	"github.com/padok-team/git-secret-scanner/internal/scan/scanners/gitleaks"
	"github.com/padok-team/git-secret-scanner/internal/scan/scanners/trufflehog"
	"github.com/padok-team/git-secret-scanner/internal/scm"
	"github.com/padok-team/git-secret-scanner/internal/scm/git"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	scanArgs  scan.ScanArgs
	scmConfig scm.ScmConfig

	noArchived  bool
	sshClone    bool
	filesOnly   bool
	noProgress  bool
	noBareClone bool

	verbose bool
	noColor bool

	scannerGitleaksPath   string
	scannerTrufflehogPath string
)

func preRun(cmd *cobra.Command, args []string) {
	logging.SetupLogger(verbose, noColor)

	// check that gitleaks and trufflehog binary exists
	gPath, err := exec.LookPath(scannerGitleaksPath)
	if err != nil {
		log.Fatal().Msgf("gitleaks binary not found: %v", err)
	}
	tPath, err := exec.LookPath(scannerTrufflehogPath)
	if err != nil {
		log.Fatal().Msgf("trufflehog binary not found: %v", err)
	}

	gitleaks.SetCommandPath(gPath)
	trufflehog.SetCommandPath(tPath)

	ok, gVersion, err := gitleaks.IsVersionValid()
	if err != nil {
		log.Warn().Msgf("failed to read gitleaks version: %v", err)
	} else if !ok {
		log.Warn().
			Msgf("this tool is designed to run with gitleaks %s or later, found gitleaks %s", gitleaks.MinVersion, gVersion)
	}
	ok, tVersion, err := trufflehog.IsVersionValid()
	if err != nil {
		log.Warn().Msgf("failed to read trufflehog version: %v", err)
	} else if !ok {
		log.Warn().
			Msgf("this tool is designed to run with trufflehog %s or later, found trufflehog %s", trufflehog.MinVersion, tVersion)
	}

	log.Debug().
		Str("gitleaks_path", gPath).
		Str("gitleaks_version", gVersion).
		Str("trufflehog_path", tPath).
		Str("trufflehog_version", tVersion).
		Msgf("running with gitleaks %s and trufflehog %s", gVersion, tVersion)

	// parse flags
	scmConfig.IncludeArchived = !noArchived
	if sshClone {
		scmConfig.GitProtocol = git.ProtocolSSH
	}
	if filesOnly {
		scanArgs.ScanType = scan.ScanTypeFilesOnly
	}
	scanArgs.ShowProgress = !noProgress
	scanArgs.BareClone = !noBareClone

	log.Debug().
		Str("scan_type", scanArgs.ScanType.String()).
		Str("report_path", scanArgs.ReportPath).
		Str("format", scanArgs.ReportFormat.String()).
		Str("fingerprints_ignore_path", scanArgs.FingerprintsIgnorePath).
		Str("baseline_path", scanArgs.BaselinePath).
		Int("max_concurrency", scanArgs.MaxConcurrency).
		Bool("show_progress", scanArgs.ShowProgress).
		Bool("bare_clone", scanArgs.BareClone).
		Msg("parsed scan args")
}

func registerCommonFlags(cmd *cobra.Command) {
	// scm config (org is not provided here as the flag name changes between scm)
	cmd.Flags().VarP(&scmConfig.Visiblity, "visibility", "v", "Visibility of repositories to scan")
	cmd.Flags().StringVarP(&scmConfig.Server, "server", "s", "", "Hostname of the server")
	cmd.Flags().BoolVar(&noArchived, "no-archived", false, "Skip archived repositories")
	cmd.Flags().BoolVar(&sshClone, "ssh-clone", false, "Use SSH to clone repositories instead of HTTPS")

	// scan args
	cmd.Flags().StringVarP(&scanArgs.ReportPath, "report-path", "r", "", "Path to the CSV report file to generate (default \"report.{json,csv}\")")
	cmd.Flags().VarP(&scanArgs.ReportFormat, "format", "f", "Format of the report")
	cmd.Flags().StringVarP(&scanArgs.FingerprintsIgnorePath, "fingerprints-ignore-path", "i", "", "Path to file with newline separated fingerprints (SHA-256) of secrets to ignore during the scan")
	cmd.Flags().StringVarP(&scanArgs.BaselinePath, "baseline-path", "b", "", "Path to the CSV report to use as baseline for the scan")
	cmd.Flags().IntVar(&scanArgs.MaxConcurrency, "max-concurrency", 5, "Maximum number of concurrent workers")
	cmd.Flags().StringVar(&scannerGitleaksPath, "scanner-gitleaks-path", "gitleaks", "Custom path to the gitleaks binary")
	cmd.Flags().StringVar(&scannerTrufflehogPath, "scanner-trufflehog-path", "trufflehog", "Custom path to the trufflehog binary")
	cmd.Flags().BoolVar(&filesOnly, "files-only", false, "Only run the scan on the files of the default branch")
	cmd.Flags().BoolVar(&noProgress, "no-progress", false, "Hide progress bar during scan")
	cmd.Flags().BoolVar(&noBareClone, "no-bare-clone", false, "Clone repositories with working directory")

	// log flags
	cmd.Flags().BoolVar(&noColor, "no-color", false, "Disable color output")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Show verbose output")

	// help flag
	cmd.Flags().BoolP("help", "h", false, fmt.Sprintf("Help for %s", cmd.Name()))
}
