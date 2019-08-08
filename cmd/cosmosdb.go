package cmd

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/DharmaOfCode/goaz/utils"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type CosmosDbState struct {
	AccountName    string
	QueueName      string
	Key            string
	All            bool
	Stalk          bool
}

var (
	comosdbState AuditSate

	cosmosdbCmd = &cobra.Command{
		Use:   "cosmosdb",
		Short: "cosmosdb",
		Long:  `audit cosmosdb`,
		Args: func(cmd *cobra.Command, args []string) error {
			if SubscriptionId == "" {
				return errors.New("please specify the subscription ID")
			}
			return nil
		},
		Run: func (cmd *cobra.Command, args []string) {
			if comosdbState.All {
				state := CosmosDbState{}
				state.Audit()
			}
		},
	}
)

func init() {
	auditState = AuditSate{}
	cosmosdbCmd.Flags().StringVar(&comosdbState.AccountName, "account", "", "Account Name")
	cosmosdbCmd.Flags().StringVar(&comosdbState.QueueName, "queue", "", "Queue Name")
	cosmosdbCmd.Flags().StringVar(&comosdbState.Key, "key", "", "Primary key for queue")
	cosmosdbCmd.Flags().BoolVarP(&comosdbState.All, "Audit all storage options", "A", false, "-A")
	cosmosdbCmd.Flags().BoolVarP(&comosdbState.Stalk, "Stalk queue", "S", false, "-S")

	rootCmd.AddCommand(cosmosdbCmd)
}


func (CosmosDbState) Audit() {
	groupsClient := resources.NewGroupsClient(SubscriptionId)
	sgClient := network.NewSecurityGroupsClient(SubscriptionId)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Println(err)
	}

	groupsClient.Authorizer = authorizer
	sgClient.Authorizer = authorizer

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Get the groups first, so we can match groups and storage accounts
	groupIterator, err := groupsClient.ListComplete(ctx, "", nil)
	var groups []string
	for list := groupIterator; list.NotDone(); err = list.NextWithContext(ctx) {
		if err != nil {
			log.Fatalf("got error: %s\n", err)
		}
		rgName := *list.Value().Name
		groups = append(groups, rgName)
	}


	columns := []string{"FROM ADDRESS", "FROM PORT", "TO ADDRESS", "TO PORT"}
	resultTable := printer.ResultTable{
		Columns: columns,
	}
	printer.Data("********Security Groups********\n")
	for _, rgName := range groups {
		sgIterator, err := sgClient.ListComplete(ctx, rgName)
		if err != nil{
			log.Println(err)
		}
		for list := sgIterator; list.NotDone(); err = list.NextWithContext(ctx) {
			if err != nil {
				log.Fatalf("got error: %s\n", err)
			}
			sgName := *list.Value().Name
			printer.InfoHeading("\t- Security Group Name: %s\n", sgName)
			rules := *list.Value().SecurityRules
			for _, rule := range rules{
				if rule.DestinationPortRange != nil{ // TODO: Check with a better propertys
					row := []string{*rule.SourceAddressPrefix, *rule.SourcePortRange, *rule.DestinationAddressPrefix, *rule.DestinationPortRange}
					resultTable.Rows = append(resultTable.Rows, row)
				}
			}
			printer.PrintTable(false, &resultTable)
		}
	}
}
