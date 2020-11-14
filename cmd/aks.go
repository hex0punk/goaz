package cmd

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hex0punk/goaz/api"
	"github.com/hex0punk/goaz/utils"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type AksState struct {
}

var (
	aksState AksState

	aksCmd = &cobra.Command{
		Use:   "aks",
		Short: "aks",
		Long:  `aks`,
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func (cmd *cobra.Command, args []string) {
			aksState.Audit()
		},
	}
)

func init() {
	rootCmd.AddCommand(aksCmd)
}


func (s *AksState) Audit() {
	groupsClient := resources.NewGroupsClient(SubscriptionId)
	aksClient := containerservice.NewManagedClustersClient(SubscriptionId)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Println(err)
	}

	groupsClient.Authorizer = authorizer
	aksClient.Authorizer = authorizer

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Get the groups first
	groups := api.GetResourceGroups(authorizer, SubscriptionId)

	columns := []string{"NAME","FQDN"}
	printer.Data("********AKS********\n")

	var aksList []containerservice.ManagedCluster
	for _, rgName := range groups {
		aksIterator, err := aksClient.ListByResourceGroup(ctx, rgName)
		if err != nil{
			log.Println(err)
		}

		for list := aksIterator; list.NotDone(); err = list.NextWithContext(ctx) {
			if err != nil {
				log.Fatalf("got error: %s\n", err)
			}
			aksList = append(aksList, list.Values()...)
		}
	}

	resultTable := printer.ResultTable{
		Columns: columns,
	}
	// TODO: evaluate whether VNET and FW rules are in effect
	for _, aks := range aksList {
		name := *aks.Name
		fqdn := *aks.Fqdn
		row := []string{
			name,
			fqdn,
		}
		resultTable.Rows = append(resultTable.Rows, row)
	}
	if len(resultTable.Rows) > 0 {
		printer.PrintTable(&resultTable)
	}
}
