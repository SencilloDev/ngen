package cmd

import (
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:               "generate",
	Short:             "Generates generate diagrams or specs",
	PersistentPreRunE: bindGenerateCmdFlags,
}

var requiredGenerateFlags = []string{
	"name",
	"id",
}

func init() {
	rootCmd.AddCommand(generateCmd)
	natsFlags(generateCmd)
	generateFlags(generateCmd)
}

func bindGenerateCmdFlags(cmd *cobra.Command, args []string) error {
	bindGenerateFlags(cmd)
	bindNatsFlags(cmd)

	return validateEnvs(requiredGenerateFlags...)
}
