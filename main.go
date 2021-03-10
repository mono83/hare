package main

import (
	"github.com/mono83/hare/cmd"
	"github.com/mono83/xray/std/xcobra"
)

func main() {
	xcobra.Start(cmd.MainCmd)
}
