package cmd

import (
	"github.com/mono83/hare/mapping"
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
	"os"
)

var downloadDelete bool
var downloadAppend bool

var downloadCmd = &cobra.Command{
	Use:     "download queue filename",
	Aliases: []string{"save", "dump", "down", "flush"},
	Short:   "Saves all messages from queue to local file",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, a []string) (err error) {
		queue := a[0]
		var f *os.File
		f, err = openFile(a[1], downloadAppend)
		if err != nil {
			return err
		}
		defer f.Close()

		return withChannel(func(ch *amqp.Channel) error {
			count := 0
			var messages []mapping.Message
			defer restore(ch, queue, &messages)

			for {
				d, ok, err := ch.Get(queue, true)
				if err != nil {
					return err
				}
				if !ok {
					xray.BOOT.Info("No more messages")
					break
				}
				m := mapping.FromDelivery(d)
				if !downloadDelete {
					messages = append(messages, m)
				}
				err = m.Fprint(f, false)
				if err != nil {
					return err
				}
				count++
				if count%100 == 0 {
					xray.BOOT.Info("Already downloaded :count messages", args.Count(count))
				}
			}
			xray.BOOT.Info("Downloaded :count messages in total", args.Count(count))
			return nil
		})
	},
}

func init() {
	downloadCmd.Flags().BoolVarP(
		&downloadDelete,
		"delete",
		"d",
		false,
		"If true, will delete messages from queue",
	)
	downloadCmd.Flags().BoolVarP(
		&downloadAppend,
		"append",
		"a",
		false,
		"If true will append to file instead of replace it",
	)
}
