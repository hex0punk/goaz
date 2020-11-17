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
		Long:  `vm`,
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


	columns := []string{"NAME","RESOURCE GROUP", "HOST ENCRYPTION", "BOOT DIAGNOSTICS", "SECURITY GROUPS"}
	resultTable := printer.ResultTable{
		Columns: columns,
	}
	printer.Data("********VM Scale Sets********\n")

	for _, rgName := range groups {
		vmScaleSetsList, err := computeClient.List(ctx, rgName)
		if err != nil{
			log.Println(err)
		}
		for _, vmss := range vmScaleSetsList.Values() {
			name := *vmss.Name
			vmProfile := vmss.VirtualMachineScaleSetProperties.VirtualMachineProfile
			encryptionAtHost := "!false"
			bootDiagnostics := "!false"
			if vmProfile.SecurityProfile != nil && *vmProfile.SecurityProfile.EncryptionAtHost {
				encryptionAtHost = "true"
			}
			if vmProfile.SecurityProfile != nil && !*vmProfile.DiagnosticsProfile.BootDiagnostics.Enabled {
				bootDiagnostics = "true"
			}
			securityGroups := "!None"
			sgLengh := len(*vmProfile.NetworkProfile.NetworkInterfaceConfigurations)
			if sgLengh > 0 {
				securityGroups = ""
				for i, sg := range *vmProfile.NetworkProfile.NetworkInterfaceConfigurations {
					if i < sgLengh - 1 {
						securityGroups = securityGroups + *sg.Name + ","
					} else {
						securityGroups = securityGroups + *sg.Name
					}
				}
			}

			row := []string{
				name,
				rgName,
				encryptionAtHost,
				bootDiagnostics,
				securityGroups,
			}
			resultTable.Rows = append(resultTable.Rows, row)
		}
	}
	if len(resultTable.Rows) > 0 {
		printer.PrintTable(&resultTable, PrintMarkdown)
	}
}