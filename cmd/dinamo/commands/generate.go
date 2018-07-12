package commands

import (
	"github.com/spf13/cobra"

	"github.com/kenjones-cisco/dinamo/generator"
)

type genOptions struct {
	Template string
	File     string
	DataFile string
	UseEnv   bool
}

func newCommandGenerate() *cobra.Command {
	opts := &genOptions{}

	action := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Generate a file",
		Long:    "Generate a file",
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sources := &generator.DataSources{
				Data:     args,
				DataFile: opts.DataFile,
				UseEnv:   opts.UseEnv,
			}
			if err := generator.Generate(opts.Template, opts.File, sources); err != nil {
				return err
			}
			cmd.Printf("generated: %s\n", opts.File)
			return nil
		},
	}

	action.Flags().StringVarP(&opts.Template, "template", "t", "", "Template file path")
	action.Flags().StringVarP(&opts.File, "file", "f", "", "Path to generated file")
	action.Flags().StringVarP(&opts.DataFile, "data", "d", "", "Path to data file")
	action.Flags().BoolVarP(&opts.UseEnv, "env", "e", false, "Use environment variables for placeholders")

	_ = action.MarkFlagRequired("template")
	_ = action.MarkFlagRequired("file")

	return action
}
