package cmd

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/hex0punk/goaz/api"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type StalkState struct {
	AccountName string
	QueueName   string
	Key         string

	Queue bool
}

var (
	stalkState StalkState

	stalkCmd = &cobra.Command{
		Use:   "stalk",
		Short: "stalk storage",
		Long:  `stalk storage`,
		Args: func(cmd *cobra.Command, args []string) error {
			if SubscriptionId == "" {
				return errors.New("please specify the subscription ID")
			}
			return nil
		},
		Run: stalk,
	}
)

func init() {
	stalkState = StalkState{}
	stalkCmd.Flags().StringVar(&stalkState.AccountName, "account", "", "Account Name")
	stalkCmd.Flags().StringVar(&stalkState.QueueName, "name", "", "Queue Name")
	stalkCmd.Flags().StringVar(&stalkState.Key, "key", "", "Primary key for queue")
	stalkCmd.Flags().BoolVarP(&stalkState.Queue, "queue", "q", false, "Stalk a queue")

	rootCmd.AddCommand(stalkCmd)
}

func stalk(cmd *cobra.Command, args []string) {
	if stalkState.Queue {
		stalkQueue()
	}
}

func stalkQueue() {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	lastMessage := ""
	for {
		fmt.Printf(".")
		messages := api.PeekMessages(ctx, stalkState.AccountName, stalkState.QueueName, stalkState.Key, 1)
		if messages.NumMessages() > 0 {
			data, err := base64.StdEncoding.DecodeString(messages.Message(int32(0)).Text)
			if err != nil {
				log.Println("[+] ", err.Error())
			}
			if string(data) != lastMessage {
				fmt.Printf("\nFOUND MESSAGE: \n")
				fmt.Println("===============================================================================\n\n\n\n\n\n")
				fmt.Print(string(data))
				fmt.Println("\n\n\n\n\n\n===============================================================================")
			}
			lastMessage = string(data)
		}
	}
}
