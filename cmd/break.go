package cmd

import (
	"github.com/mono83/hare/mapping"
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
	"os"
)

var breakAppend bool
var breakPriority int

var breakCmd = &cobra.Command{
	Use:     "circuit-breaker queue filename",
	Aliases: []string{"cb", "break", "circuit-break"},
	Args:    cobra.ExactArgs(2),
	Short:   "Runs circuit breaker, that will intercept all messages from queue and flush them to file",
	RunE: func(cmd *cobra.Command, a []string) (err error) {
		queue := a[0]
		var f *os.File
		f, err = openFile(a[1], breakAppend)
		if err != nil {
			return err
		}
		defer f.Close()

		return withChannel(func(ch *amqp.Channel) error {
			xray.BOOT.Info(
				"Starting circuit breaker with priority :count on queue :name",
				args.Name(queue),
				args.Count(breakPriority),
			)
			delivery, err := ch.Consume(queue, "", true, false, false, false, amqp.Table{
				"x-priority": breakPriority,
			})
			if err != nil {
				return err
			}
			for d := range delivery {
				m := mapping.FromDelivery(d)
				err = m.Fprint(f, false)
				if err != nil {
					return err
				}
			}
			return nil
		})
	},
}

func init() {
	breakCmd.Flags().IntVarP(
		&breakPriority,
		"priority",
		"p",
		10,
		"Consumer priority",
	)
	breakCmd.Flags().BoolVarP(
		&breakAppend,
		"append",
		"a",
		false,
		"If true will append to file instead of replace it",
	)
}
