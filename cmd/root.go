package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	Verbose        bool
	SubscriptionId string
	PrintMarkdown  bool
	rootCmd        = &cobra.Command{
		Use:   "goaz",
		Short: "Azure security auditor and stalker",
		Long:  `Something or other Azure something something`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must enter at least one arg")
			}
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&SubscriptionId, "subscriptionId", "", "subscription ID to use")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&PrintMarkdown, "markdown", "m", false, "Print markdown format")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
