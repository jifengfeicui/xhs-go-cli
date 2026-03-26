package main

import (
	"os"

	"xhs-go-cli/internal/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	rootCmd.AddCommand(
		cmd.NewImportSourcesCmd(),
		cmd.NewQueryGenCmd(),
		cmd.NewSearchCmd(),
		cmd.NewFetchDetailCmd(),
		cmd.NewQualifyCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
