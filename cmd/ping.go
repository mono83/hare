package cmd

import (
	"github.com/mono83/xray"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:     "ping",
	Aliases: []string{"test"},
	Short:   "Tests RabbitMQ connection",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := dial()
		if err == nil {
			err = conn.Close()
		}
		if err == nil {
			xray.BOOT.Info("RabbitMQ connection is OK")
		}
		return err
	},
}
