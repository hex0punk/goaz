package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type NetState struct {
	AuditAll            bool
}

var (
	netState NetState

	netCmd = &cobra.Command{
		Use:   "net",
		Short: "net",
		Long:  `net`,
		Args: func(cmd *cobra.Command, args []string) error {
			if SubscriptionId == "" {
				return errors.New("please specify the subscription ID")
			}
			return nil
		},
		Run: func (cmd *cobra.Command, args []string) {
			fmt.Println("test")
		},
	}
)

func init() {
	rootCmd.AddCommand(netCmd)
}
