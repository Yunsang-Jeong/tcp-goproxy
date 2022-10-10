package main

import "tcp-proxy/cmd"

func main() {
	start := &cmd.StartCmd{}

	cmd.RootCmd.AddCommand(start.Init())
	cmd.Execute()
}
