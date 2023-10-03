package cmd

import (
	"github.com/mono83/hare/mapping"
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
	"os"
)

var rabbitMQDSN string

// MainCmd is main entry point
var MainCmd = &cobra.Command{
	Use:   "hare",
	Short: "RabbitMQ toolkit",
}

func init() {
	MainCmd.PersistentFlags().StringVarP(
		&rabbitMQDSN,
		"uri",
		"u",
		"amqp://guest:guest@localhost:5672/",
		"RabbitMQ URI",
	)

	MainCmd.AddCommand(
		pingCmd,
		visitCmd,
		copyCmd,
		moveCmd,
		downloadCmd,
		uploadCmd,
		breakCmd,
	)
}

func dial() (*amqp.Connection, error) {
	return amqp.Dial(rabbitMQDSN)
}

func withChannel(f func(channel *amqp.Channel) error) error {
	conn, err := dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return f(ch)
}

func restore(ch *amqp.Channel, queue string, messages *[]mapping.Message) {
	if messages != nil {
		xray.BOOT.Info("Restoring :count messages in :name", args.Count(len(*messages)), args.Name(queue))
		for _, m := range *messages {
			if err := m.Publish(ch, "", queue); err != nil {
				xray.BOOT.Error("Error restoring message :err", args.Error{Err: err})
			}
		}
	}
}

func openFile(name string, append bool) (f *os.File, err error) {
	if append {
		xray.BOOT.Info("Appending to :name", args.Name(name))
		f, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		xray.BOOT.Info("Creating/replacing :name", args.Name(name))
		f, err = os.Create(name)
	}
	return
}
