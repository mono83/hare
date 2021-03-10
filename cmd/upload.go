package cmd

import (
	"bufio"
	"github.com/mono83/hare/mapping"
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
	"os"
	"strings"
)

var uploadReplicate int
var uploadExchange string

var uploadCmd = &cobra.Command{
	Use:     "upload queue filename",
	Aliases: []string{"restore", "up"},
	Short:   "Places previously stored messages to queue",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, a []string) error {
		queue := a[0]
		f, err := os.Open(a[1])
		if err != nil {
			return err
		}
		defer f.Close()

		return withChannel(func(ch *amqp.Channel) error {
			scan := bufio.NewScanner(f)
			count := 0
			for scan.Scan() {
				line := strings.TrimSpace(scan.Text())
				if len(line) == 0 {
					continue
				}
				m, err := mapping.FromJSON(line)
				if err != nil {
					return err
				}
				for i := 0; i < uploadReplicate; i++ {
					if len(uploadExchange) > 0 {
						// Custom exchange
						if err := m.PublishCustom(ch, uploadExchange, queue); err != nil {
							return err
						}
					} else {
						if err := m.Publish(ch, queue); err != nil {
							return err
						}
					}
					count++
					if count%100 == 0 {
						xray.BOOT.Info("Already uploaded :count messages", args.Count(count))
					}
				}
			}
			xray.BOOT.Info("Uploaded :count messages in total", args.Count(count))
			return nil
		})
	},
}

func init() {
	uploadCmd.Flags().IntVarP(
		&uploadReplicate,
		"replicate",
		"r",
		1,
		"Amount of copies per single line",
	)
	uploadCmd.Flags().StringVarP(
		&uploadExchange,
		"exchange",
		"e",
		"",
		"Exchange to use for uploaded messages",
	)
}
