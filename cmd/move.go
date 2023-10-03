package cmd

import (
	"errors"
	"github.com/mono83/hare/mapping"
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
)

var moveExchange string

var moveCmd = &cobra.Command{
	Use:   "move source target",
	Short: "Moves all messages from source queue to target queue",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, a []string) error {
		source := a[0]
		target := a[1]

		if source == target {
			return errors.New("unable to move messages into same queue")
		}

		return withChannel(func(ch *amqp.Channel) error {
			count := 0
			for {
				d, ok, err := ch.Get(source, true)
				if err != nil {
					return err
				}
				if !ok {
					xray.BOOT.Info("No more messages")
					break
				}
				if err := mapping.FromDelivery(d).Publish(ch, moveExchange, target); err != nil {
					return err
				}
				count++
			}
			xray.BOOT.Info("Moved :count messages", args.Count(count))
			return nil
		})
	},
}

func init() {
	moveCmd.Flags().StringVarP(
		&moveExchange,
		"exchange",
		"e",
		"",
		"Exchange to use for moved messages",
	)
}
