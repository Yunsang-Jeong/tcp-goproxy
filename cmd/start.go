package cmd

import (
	"tcp-proxy/pkg/handler"

	"github.com/spf13/cobra"
)

type StartCmd struct{}

type startFlag struct {
	shorten      string
	defaultValue string
	description  string
	requirement  bool
}

var startFlags = map[string]startFlag{
	"Type": {
		shorten:      "t",
		defaultValue: "ftp",
		description:  "[opt] The type of proxy",
	},
}

func (s *StartCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "start",
		Short: "Start proxy server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			proxyType, _ := cmd.Flags().GetString("Type")

			switch proxyType {
			case "ftp":
				if err := handler.StartTCPHandler(":8080"); err != nil {
					return err
				}
			default:
				return nil
			}

			return nil
		},
	}

	for name, flag := range startFlags {
		c.Flags().StringP(
			name,
			flag.shorten,
			flag.defaultValue,
			flag.description,
		)

		if flag.requirement {
			c.MarkFlagRequired(name)
		}
	}

	return c
}
