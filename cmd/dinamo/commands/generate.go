package commands

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kenjones-cisco/dinamo/generator"
)

type genOptions struct {
	Template string
	File     string
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
			var err error

			if opts.UseEnv {
				err = generator.GenerateUsingEnv(opts.Template, opts.File)
			} else {
				err = generator.Generate(opts.Template, opts.File, argsMap(args))
			}

			if err == nil {
				cmd.Printf("generated: %s\n", opts.File)
			}
			return err
		},
	}

	action.Flags().StringVarP(&opts.Template, "template", "t", "", "Template file path")
	action.Flags().StringVarP(&opts.File, "file", "f", "", "Path to generated file")
	action.Flags().BoolVarP(&opts.UseEnv, "env", "e", false, "Use environment variables for placeholders")

	_ = action.MarkFlagRequired("template")
	_ = action.MarkFlagRequired("file")

	return action
}

func argsMap(args []string) map[string]interface{} {
	amap := make(map[string]interface{})
	for _, item := range args {
		kv := strings.SplitN(item, "=", 2)
		amap[kv[0]] = kv[1]
	}
	log.Debugf("generate args: %v", amap)
	return amap
}
