package cmd

import (
	"errors"
	"github.com/mono83/hare/mapping"
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

var copyCmd = &cobra.Command{
	Use:   "copy source target",
	Short: "Copies all messages from source queue to target queue",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, a []string) error {
		source := a[0]
		target := a[1]

		if source == target {
			return errors.New("unable to copy messages into same queue")
		}

		return withChannel(func(ch *amqp.Channel) error {
			count := 0
			var messages []mapping.Message
			defer restore(ch, source, &messages)
			for {
				d, ok, err := ch.Get(source, true)
				if err != nil {
					return err
				}
				if !ok {
					xray.BOOT.Info("No more messages")
					break
				}
				m := mapping.FromDelivery(d)
				messages = append(messages, m)
				if err := m.Publish(ch, target); err != nil {
					return err
				}
				count++
			}
			xray.BOOT.Info("Copied :count messages", args.Count(count))
			return nil
		})
	},
}
