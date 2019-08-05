package cmd

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/DharmaOfCode/goaz/api"
	"github.com/DharmaOfCode/goaz/utils"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type AuditSate struct {
	AccountName    string
	QueueName      string
	Key            string
	All            bool
	Stalk          bool
}

var (
	auditState AuditSate

	queryCmd = &cobra.Command{
		Use:   "audit",
		Short: "audit storage",
		Long:  `audit storage`,
		Args: func(cmd *cobra.Command, args []string) error {
			if SubscriptionId == "" {
				return errors.New("please specify the subscription ID")
			}
			return nil
		},
		Run: runAudit,
	}
)

func init() {
	auditState = AuditSate{}
	queryCmd.Flags().StringVar(&auditState.AccountName, "account", "", "Account Name")
	queryCmd.Flags().StringVar(&auditState.QueueName, "queue", "", "Queue Name")
	queryCmd.Flags().StringVar(&auditState.Key, "key", "", "Primary key for queue")
	queryCmd.Flags().BoolVarP(&auditState.All, "Audit all storage options", "A", false, "-A")
	queryCmd.Flags().BoolVarP(&auditState.Stalk, "Stalk queue", "S", false, "-S")

	rootCmd.AddCommand(queryCmd)
}

func runAudit(cmd *cobra.Command, args []string) {
	if auditState.All {
		AuditAll()
	}
}

func AuditAll() {
	groupsClient := resources.NewGroupsClient(SubscriptionId)
	storageAccountsClient := storage.NewAccountsClient(SubscriptionId)
	blobStorageClient := storage.NewBlobContainersClient(SubscriptionId)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Println(err)
	}

	storageAccountsClient.Authorizer = authorizer
	groupsClient.Authorizer = authorizer
	blobStorageClient.Authorizer = authorizer

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Get the groups first, so we can match groups and storage accounts
	groupIterator, err := groupsClient.ListComplete(ctx, "", nil)
	var groups []string
	for list := groupIterator; list.NotDone(); err = list.NextWithContext(ctx) {
		if err != nil {
			log.Fatalf("got error: %s", err)
		}
		rgName := *list.Value().Name
		groups = append(groups, rgName)
	}

	// Get storage accounts per resource group
	for _, rgName := range groups {
		storageAccounts, err := storageAccountsClient.ListByResourceGroup(ctx, rgName)
		if err != nil {
			log.Println(err)
		}

		if len(*storageAccounts.Value) == 0{
			continue
		}
		fmt.Printf("************ Storage Accounts for Group %s ************\n\n", rgName)
		for _, acc := range *storageAccounts.Value {
			printer.Info("Storage Account Data\n")
			fmt.Println("========================")
			printer.Data("- Accounts:\n")
			printer.InfoHeading("\t- Account Name: %s\n", *acc.Name)
			printer.Info("\t- ID: %s\n", *acc.ID)
			printer.Info("\t- Type: %s\n", *acc.Type)
			printer.Data("\t- Keys:\n")
			// Print keys
			keys, err := storageAccountsClient.ListKeys(ctx, rgName, *acc.Name)
			if err != nil {
				printer.Warning("\t\t[+] Unable to read keys\n")
			}

			var keyStrings []string
			if err == nil{
				for _, key := range *keys.Keys {
					keyStrings = append(keyStrings, *key.Value)
					printer.Info("\t\tKey name: %s\n\t\tValue: %s\n\t\tPermissions: %s\n",
						*key.KeyName,
						*key.Value,
						key.Permissions)
					fmt.Println("\t\t----------------")
				}
			}

			// Get containers per storage account
			containersIterator, err := blobStorageClient.ListComplete(ctx, rgName, *acc.Name, "", "", "")
			if err != nil {
				log.Println(err)
			}

			printer.Data("\t- Containers:\n")
			for list := containersIterator; list.NotDone(); err = list.NextWithContext(ctx) {
				if err != nil {
					log.Fatalf("got error: %s", err)
				}
				containerName := *list.Value().Name
				printer.InfoHeading("\t\t- Container Name: %s\n", containerName)
				switch list.Value().PublicAccess {
					case "blob", "container":
					{
						printer.Danger("\t\t- Public access type: %s\n", list.Value().PublicAccess)
					}
					default:{
						printer.Info("\t\t- Public access type: %s\n", list.Value().PublicAccess)
					}
				}

				// Get blobs
				if keyStrings == nil{
					continue
				}
				blobs, err := storageapi.ListBlobs(ctx, *acc.Name, containerName, keyStrings[0])
				if err != nil {
					log.Fatalf("got error: %s", err)
				}
				printer.Data("\t\t- Blobs:\n")
				for _, b := range blobs.Segment.BlobItems {
					printer.InfoHeading("\t\t\t- Blob Name: %s\n", b.Name)
					printer.Info("\t\t\t- URL: %s\n", storageapi.BlobURL(*acc.Name, containerName, b.Name))
					printer.Info("\t\t\t- Content Type: %s\n", *b.Properties.ContentType)
					printer.Info("\t\t\t- Type: %s\n", b.Properties.BlobType)
					fmt.Println()
					//TODO: create a url.URL and download blob, or keep it as a command line option
				}
			}

			if keyStrings == nil{
				continue
			}
			// TODO: Get tables in container

			//Get file shares.
			shareDirectoriesResponse := storageapi.GetShareDirectories(ctx, *acc.Name, keyStrings[0])
			printer.Data("\t- Shares:\n")
			for _, d := range shareDirectoriesResponse.ShareItems {
				printer.InfoHeading("\t\t- Share Name: %s\n", d.Name)
				shares, err := storageapi.ListFiles(ctx, *acc.Name, d.Name, keyStrings[0])
				if err != nil {
					log.Println(err)
				}
				printer.Data("\t\t- Files:\n")
				for _, s := range shares.FileItems {
					printer.InfoHeading("\t\t\t-File Name: %s\n", s.Name)
					printer.Info("\t\t\t-URL: %s\n", storageapi.FileURL(*acc.Name, d.Name, s.Name))
				}
				fmt.Println()
			}

			// Get Queues
			printer.Data("\t- Queues:\n")
			queuesResponse := storageapi.ListQueues(ctx, *acc.Name, keyStrings[0])
			for _, q := range queuesResponse.QueueItems {
				printer.InfoHeading("\t\t- Queue Name: %s\n", q.Name)
				messages := storageapi.PeekMessages(ctx, *acc.Name, q.Name, keyStrings[0], 32)
				if err != nil {
					log.Println(err)
				}

				printer.Data("\t\t- Messages:\n")
				for i := 0; i < int(messages.NumMessages()); i++ {
					data, err := base64.StdEncoding.DecodeString(messages.Message(int32(i)).Text)
					if err != nil {
						log.Println("[+] ", err.Error())
					}
					printer.Info("\t\t\t- %s\n", string(data))
				}
				fmt.Println()
			}

			fmt.Print("\n\n")
		}
	}
}
