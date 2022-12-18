package cmd

import (
	"github.com/spf13/cobra"
)

type StartCmd struct{}

func (s *StartCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "start",
		Short: "Start proxy server.",
	}

	startTCP := &StartTCPCmd{}
	startReadOnlyFTP := &StartReadOnlyFTPCmd{}

	c.AddCommand(startTCP.Init())
	c.AddCommand(startReadOnlyFTP.Init())

	return c
}
