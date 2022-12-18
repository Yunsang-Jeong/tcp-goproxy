package cmd

import (
	"github.com/Yunsang-Jeong/tcp-goproxy/internal/proxy"

	"github.com/spf13/cobra"
)

type StartReadOnlyFTPCmd struct{}

var startReadOnlyFTPFlags = map[string]flag{}

func (s *StartReadOnlyFTPCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "read-only-ftp",
		Short: "Start read-only ftp proxy server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ps := proxy.NewFTPProxyServer(":8080")
			if err := ps.Start(); err != nil {
				return err
			}

			return nil
		},
	}
	cobraFlagRegister(c, startReadOnlyFTPFlags)

	return c
}
