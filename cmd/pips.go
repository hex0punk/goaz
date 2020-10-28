package cmd

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hex0punk/goaz/utils"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type PipsState struct {
}

var (
	pipsState PipsState

	pipsCmd = &cobra.Command{
		Use:   "pips",
		Short: "pips",
		Long:  `pips`,
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func (cmd *cobra.Command, args []string) {
			pipsState.Audit()
		},
	}
)

func init() {
	netCmd.AddCommand(pipsCmd)
}


func (s *PipsState) Audit() {
	groupsClient := resources.NewGroupsClient(SubscriptionId)
	netClient := network.NewPublicIPAddressesClient(SubscriptionId)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Println(err)
	}

	groupsClient.Authorizer = authorizer
	netClient.Authorizer = authorizer

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Get the groups first
	groupIterator, err := groupsClient.ListComplete(ctx, "", nil)
	var groups []string
	for list := groupIterator; list.NotDone(); err = list.NextWithContext(ctx) {
		if err != nil {
			log.Fatalf("got error: %s\n", err)
		}
		rgName := *list.Value().Name
		groups = append(groups, rgName)
	}


	columns := []string{"NAME","PUBLIC IP","LOCATION"}
	printer.Data("********Public IPs********\n")

	var publicIPList []network.PublicIPAddress
	for _, rgName := range groups {
		sgIterator, err := netClient.ListComplete(ctx, rgName)
		if err != nil{
			log.Println(err)
		}

		for list := sgIterator; list.NotDone(); err = list.NextWithContext(ctx) {
			if err != nil {
				log.Fatalf("got error: %s\n", err)
			}
			publicIPList = append(publicIPList, list.Value())
		}
	}

	resultTable := printer.ResultTable{
		Columns: columns,
	}
	for _, ip := range publicIPList {
		ipName := *ip.Name
		publicAddress := *ip.IPAddress
		location := *ip.Location

		row := []string{
			ipName,
			publicAddress,
			location,
		}
		resultTable.Rows = append(resultTable.Rows, row)
	}
	if len(resultTable.Rows) > 0 {
		printer.PrintTable(&resultTable)
	}
}
