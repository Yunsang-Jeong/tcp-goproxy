package main

import "github.com/Yunsang-Jeong/tcp-goproxy/cmd"

func main() {
	start := &cmd.StartCmd{}

	cmd.RootCmd.AddCommand(start.Init())
	cmd.Execute()
}
