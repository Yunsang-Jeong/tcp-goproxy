package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:           "tcproxy",
	Short:         "TCP-Proxy",
	Long:          "TCP-Proxy",
	SilenceErrors: true,
}

func Execute() {
	// Remove help for root command
	RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Remove shell completion
	RootCmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
		HiddenDefaultCmd:    true,
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}