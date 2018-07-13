package commands

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kenjones-cisco/dinamo/version"
)

type commonOptions struct {
	Debug    bool
	LogLevel string
	Version  bool
}

// NewCommandCLI creates the root command
func NewCommandCLI() *cobra.Command {
	opts := &commonOptions{}

	rootCmd := &cobra.Command{
		Use:          version.ShortName,
		Short:        version.ProductName,
		Long:         "Lightweight command-line utility for generating file(s) from using go templates.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Version {
				cmd.Print(version.GetVersionDisplay())
				return nil
			}
			cmd.Println("")
			cmd.Println(cmd.UsageString())

			return nil
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			var loglevel logrus.Level

			if opts.LogLevel != "" {
				loglevel, err = logrus.ParseLevel(strings.ToLower(opts.LogLevel))
				if err != nil {
					cmd.Println("Unknown log-level provided:", opts.LogLevel)
					loglevel = logrus.InfoLevel
				}
			}

			if opts.Debug {
				loglevel = logrus.DebugLevel
			}

			logrus.SetLevel(loglevel)
			logrus.SetFormatter(&logrus.TextFormatter{
				DisableTimestamp: false,
				FullTimestamp:    true,
				DisableSorting:   true,
			})
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&opts.Debug, "debug", "D", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&opts.LogLevel, "log-level", "l", "info", `Set the logging level ("debug", "info", "warn", "error", "fatal")`)
	rootCmd.Flags().BoolVarP(&opts.Version, "version", "v", false, "Print version information and quit")

	rootCmd.AddCommand(newCommandGenerate())

	return rootCmd
}

// Execute is the entrypoint to run any command
func Execute() {
	cmd := NewCommandCLI()
	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
