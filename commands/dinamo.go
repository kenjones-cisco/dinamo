package commands

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kenjones-cisco/dinamo/version"
)

type commonOptions struct {
	LogLevel string
	Debug    bool
	Version  bool
}

// NewCommandCLI creates the root command.
func NewCommandCLI() *cobra.Command {
	opts := &commonOptions{}

	rootCmd := &cobra.Command{
		Use:          version.ShortName,
		Short:        version.ProductName,
		Long:         "Lightweight command-line utility for generating file(s) from using go templates.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.Version {
				cmd.Print(version.GetVersionDisplay())

				return nil
			}

			cmd.Println("")
			cmd.Println(cmd.UsageString())

			return nil
		},
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			logLevel := slog.LevelInfo

			if opts.LogLevel != "" {
				var ok bool

				logLevel, ok = parseLogLevel(opts.LogLevel)
				if !ok {
					cmd.Println("Unknown log-level provided:", opts.LogLevel)

					logLevel = slog.LevelInfo
				}
			}

			if opts.Debug {
				logLevel = slog.LevelDebug
			}

			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: logLevel,
			})))
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&opts.Debug, "debug", "D", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&opts.LogLevel, "log-level", "l", "info", `Set the logging level ("debug", "info", "warn", "error", "fatal")`)
	rootCmd.Flags().BoolVarP(&opts.Version, "version", "v", false, "Print version information and quit")

	rootCmd.AddCommand(newCommandGenerate())

	return rootCmd
}

// Execute is the entrypoint to run any command.
func Execute() {
	cmd := NewCommandCLI()
	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func parseLogLevel(level string) (slog.Level, bool) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, true
	case "info":
		return slog.LevelInfo, true
	case "warn", "warning":
		return slog.LevelWarn, true
	case "error", "fatal":
		return slog.LevelError, true
	default:
		return slog.LevelInfo, false
	}
}
