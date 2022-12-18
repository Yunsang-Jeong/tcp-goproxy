package cmd

import (
	"github.com/Yunsang-Jeong/tcp-goproxy/internal/proxy"

	"github.com/spf13/cobra"
)

type StartTCPCmd struct{}

var startTCPFlags = map[string]flag{
	"Target": {
		shorten:     "t",
		description: "[req] The taget address to proxy",
		requirement: true,
	},
}

func (s *StartTCPCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "tcp",
		Short: "Start proxy server for tcp.",
		RunE: func(cmd *cobra.Command, args []string) error {
			targetAddr, _ := cmd.Flags().GetString("Target")
			ps := proxy.NewTCPProxyServer(":8080", targetAddr)
			if err := ps.Start(); err != nil {
				return err
			}

			return nil
		},
	}
	cobraFlagRegister(c, startTCPFlags)

	return c
}
