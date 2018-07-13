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
		Short:   "Generate files",
		Long: `Generate files using go templates and multiple data sources.
Use any combination of the datasources JSON/YAML data files, environment variables, and key-value pairs arguments.
`,
		Example: `
# create output.txt from config.tmpl using the key-value pairs
dinamo gen -t config.tmpl -f output.txt key1=value1 key2=value2

# create output.txt from config.tmpl using the JSON data file source.json
dinamo gen -t config.tmpl -f output.txt -d source.json

# create output.txt from config.tmpl using the YAML data file source.yaml
dinamo gen -t config.tmpl -f output.txt -d source.yaml

# create output.txt from config.tmpl using environment variables
dinamo gen -t config.tmpl -f output.txt -e

# create output.txt from config.tmpl using the key-value pairs and the YAML data file source.yml
dinamo gen -t config.tmpl -f output.txt -d source.yml key1=value1 key2=value2

# create output.txt from config.tmpl using the key-value pairs and the JSON data file source.json
dinamo gen -t config.tmpl -f output.txt -d source.json key1=value1 key2=value2

# create output.txt from config.tmpl using the key-value pairs and environment variables
dinamo gen -t config.tmpl -f output.txt -e key1=value1 key2=value2

# create output.txt from config.tmpl using the key-value pairs, environment variables, and the YAML data file source.yml
dinamo gen -t config.tmpl -f output.txt -e -d source.yml key1=value1 key2=value2
`,
		Args:              cobra.ArbitraryArgs,
		DisableAutoGenTag: true,
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
	action.Flags().StringVarP(&opts.DataFile, "data", "d", "", `Path to data file of type ("json", "yaml", "yml")`)
	action.Flags().BoolVarP(&opts.UseEnv, "env", "e", false, "Use environment variables for placeholders")

	_ = action.MarkFlagRequired("template")
	_ = action.MarkFlagRequired("file")

	_ = action.MarkFlagFilename("template")
	_ = action.MarkFlagFilename("file")
	_ = action.MarkFlagFilename("data", ".json", ".yaml", ".yml")

	return action
}
