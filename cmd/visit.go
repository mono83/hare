package cmd

import (
	"github.com/mono83/hare/mapping"
	"github.com/mono83/xray"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

var visitCount int

var visitCmd = &cobra.Command{
	Use:     "visit queue",
	Aliases: []string{"view", "look"},
	Short:   "Takes some messages from queue and then requeue them",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		queue := args[0]

		return withChannel(func(ch *amqp.Channel) error {
			var messages []mapping.Message
			defer restore(ch, queue, &messages)

			for count := visitCount; count > 0; count-- {
				delivery, ok, err := ch.Get(queue, true)
				if err != nil {
					return err
				}
				if !ok {
					xray.BOOT.Info("No more messages")
					break
				}
				m := mapping.FromDelivery(delivery)
				_ = m.Fprint(nil, true)
				messages = append(messages, m)
			}

			return nil
		})
	},
}

func init() {
	visitCmd.Flags().IntVarP(&visitCount, "count", "c", 1, "Count to get")
}
