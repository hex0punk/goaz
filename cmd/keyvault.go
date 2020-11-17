package cmd

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hex0punk/goaz/api"
	"github.com/hex0punk/goaz/utils"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"time"
)

type KeyState struct {
}

var (
	keyState KeyState

	keyCmd = &cobra.Command{
		Use:   "keyvaults",
		Short: "kv",
		Long:  `keyvaults`,
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func (cmd *cobra.Command, args []string) {
			keyState.Audit()
		},
	}
)

func init() {
	rootCmd.AddCommand(keyCmd)
}


func (s *KeyState) Audit() {
	keyClient := keyvault.NewVaultsClient(SubscriptionId)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Println(err)
	}

	keyClient.Authorizer = authorizer

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Get the groups first
	groups := api.GetResourceGroups(authorizer, SubscriptionId)

	columns := []string{"NAME","RESOURCE GROUP","URI", "FOR DEPLOYMENTS", "FOR ENCRYPTION", "VNET RESTRICTED"}
	printer.Data("********Key Vaults********\n")

	resultTable := printer.ResultTable{
		Columns: columns,
	}
	for _, rgName := range groups {
		var top int32
		top = 99999
		keys, err := keyClient.ListByResourceGroup(ctx, rgName, &top)
		if err != nil {
			log.Println(err.Error())
		}
		for _, key := range keys.Values() {
			name := *key.Name
			hasFw := "!false"
			enabledForDeployment := "false"
			if key.Properties.EnabledForDeployment != nil {
				enabledForDeployment = strconv.FormatBool(*key.Properties.EnabledForDeployment)
			}
			enabledForDiskEnc := "false"
			if key.Properties.EnabledForDiskEncryption != nil {
				enabledForDiskEnc = strconv.FormatBool(*key.Properties.EnabledForDiskEncryption)
			}
			uri := *key.Properties.VaultURI

			if key.Properties.NetworkAcls != nil {
				if len(*key.Properties.NetworkAcls.VirtualNetworkRules) > 0 || len(*key.Properties.NetworkAcls.IPRules) > 0{
					hasFw = "true"
				}
			}
			row := []string{
				name,
				rgName,
				uri,
				enabledForDeployment,
				enabledForDiskEnc,
				hasFw,
			}
			resultTable.Rows = append(resultTable.Rows, row)
		}
	}
	if len(resultTable.Rows) > 0 {
		printer.PrintTable(&resultTable, PrintMarkdown)
	}
}
