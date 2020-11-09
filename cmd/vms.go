package cmd

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-30/compute"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hex0punk/goaz/api"
	"github.com/hex0punk/goaz/utils"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type VMState struct {
	AuditAll            bool
}

var (
	vmState      VMState

	vmCmd = &cobra.Command{
		Use:   "vm",
		Short: "vm",
		Long:  `audit vm`,
		Args: func(cmd *cobra.Command, args []string) error {
			if SubscriptionId == "" {
				return errors.New("please specify the subscription ID")
			}
			return nil
		},
		Run: func (cmd *cobra.Command, args []string) {
			if vmState.AuditAll {
				vmState.Audit()
			}
		},
	}
)

func init() {
	vmCmd.Flags().BoolVarP(&vmState.AuditAll, "Audit", "A", false, "-A")

	rootCmd.AddCommand(vmCmd)
}


func (s *VMState) Audit() {
	computeClient := compute.NewVirtualMachineScaleSetsClient(SubscriptionId)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Println(err)
	}
	
	computeClient.Authorizer = authorizer

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	groups := api.GetResourceGroups(authorizer, SubscriptionId)


	columns := []string{"NAME","UUID"}
	printer.Data("********Public IPs********\n")

	var vmScaleSetsList []compute.VirtualMachineScaleSet
	for _, rgName := range groups {
		vmssIterator, err := computeClient.ListComplete(ctx, rgName)
		if err != nil{
			log.Println(err)
		}

		for list := vmssIterator; list.NotDone(); err = list.NextWithContext(ctx) {
			if err != nil {
				log.Fatalf("got error: %s\n", err)
			}
			vmScaleSetsList = append(vmScaleSetsList, list.Value())
		}
	}

	resultTable := printer.ResultTable{
		Columns: columns,
	}
	for _, vmss := range vmScaleSetsList {
		name := *vmss.Name
		row := []string{
			name,
			*vmss.VirtualMachineScaleSetProperties.UniqueID,
		}
		resultTable.Rows = append(resultTable.Rows, row)
	}
	if len(resultTable.Rows) > 0 {
		printer.PrintTable(&resultTable)
	}
}