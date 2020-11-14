package cmd

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hex0punk/goaz/api"
	"github.com/hex0punk/goaz/utils"
	"github.com/spf13/cobra"
	"log"
	"strconv"
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
	netClient := network.NewPublicIPAddressesClient(SubscriptionId)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Get the groups first
	groups := api.GetResourceGroups(authorizer, SubscriptionId)


	columns := []string{"NAME","PUBLIC IP","LOCATION", "DOMAIN", "Anti-DDoS"}
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

		domain := ""
		if ip.DNSSettings != nil {
			domain = *ip.DNSSettings.Fqdn
		}

		doSProtected := "!false"
		if ip.DdosSettings != nil {
			doSProtected = strconv.FormatBool(*ip.DdosSettings.ProtectedIP)
		}

		row := []string{
			ipName,
			publicAddress,
			location,
			domain,
			doSProtected,
		}
		resultTable.Rows = append(resultTable.Rows, row)
	}
	if len(resultTable.Rows) > 0 {
		printer.PrintTable(&resultTable, PrintMarkdown)
	}
}
