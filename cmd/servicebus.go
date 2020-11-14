package cmd

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/preview/servicebus/mgmt/2018-01-01-preview/servicebus"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hex0punk/goaz/api"
	"github.com/hex0punk/goaz/utils"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type SBusState struct {
}

var (
	sbusState SBusState

	sbusCmd = &cobra.Command{
		Use:   "sbus",
		Short: "sbus",
		Long:  `sbus`,
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func (cmd *cobra.Command, args []string) {
			sbusState.Audit()
		},
	}
)

func init() {
	rootCmd.AddCommand(sbusCmd)
}


func (s *SBusState) Audit() {
	sbNamespaceClient := servicebus.NewNamespacesClient(SubscriptionId)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Println(err)
	}

	sbNamespaceClient.Authorizer = authorizer

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Get the groups first
	groups := api.GetResourceGroups(authorizer, SubscriptionId)


	columns := []string{"NAME", "RESOURCE GROUP", "ENDPOINT","REDUNDANT", "VNET"}
	printer.Data("********Service Bus********\n")

	resultTable := printer.ResultTable{
		Columns: columns,
	}

	for _, rgName := range groups {
		sbNS, err := sbNamespaceClient.ListByResourceGroup(ctx, rgName)
		if err != nil {
			log.Println(err)
		}
		for _, sb := range sbNS.Values() {
			name := *sb.Name
			endpoint := *sb.ServiceBusEndpoint
			redundant := "!false"
			if *sb.SBNamespaceProperties.ZoneRedundant{
				redundant = "true"
			}
			//vaultUri := *sb.SBNamespaceProperties.Encryption.KeyVaultProperties.KeyVaultURI
			hasVnetRules := "!false"
			vnetRulesPage, err := sbNamespaceClient.ListVirtualNetworkRules(ctx, rgName, name)
			if err != nil {
				log.Println("[+] ", err.Error())
			}
			if &vnetRulesPage != nil {
				vnetRules := vnetRulesPage.Values()
				if len(vnetRules) > 0 {
					hasVnetRules = "true"
				}
			}
			row := []string{
				name,
				rgName,
				endpoint,
				redundant,
				hasVnetRules,
			}
			resultTable.Rows = append(resultTable.Rows, row)
		}
	}

	if len(resultTable.Rows) > 0 {
		printer.PrintTable(&resultTable, PrintMarkdown)
	}
}
