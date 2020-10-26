package cmd

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/graphrbac/graphrbac"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hex0punk/goaz/utils"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type GraphRbacState struct {
	AuditAll            bool
	TenantId		string
}

var (
	graphRbacState GraphRbacState

	graphRbacCmd = &cobra.Command{
		Use:   "graphrbac",
		Short: "graph",
		Long:  `graphd`,
		Args: func(cmd *cobra.Command, args []string) error {
			if graphRbacState.TenantId  == "" {
				return errors.New("please specify the subscription ID")
			}
			return nil
		},
		Run: func (cmd *cobra.Command, args []string) {
			if graphRbacState.AuditAll {
				graphRbacState.Audit()
			}
		},
	}
)

func init() {
	storageState = StorageState{}
	graphRbacCmd.Flags().BoolVarP(&graphRbacState.AuditAll, "Audit", "A", false, "-A")
	graphRbacCmd.PersistentFlags().StringVar(&graphRbacState.TenantId, "tenantId", "", "tenants ID to use")

	rootCmd.AddCommand(graphRbacCmd)
}


func (s *GraphRbacState) Audit() {
	grClient := graphrbac.NewApplicationsClient(graphRbacState.TenantId)
	config := auth.NewMSIConfig()
	auth.NewAuthorizerFromCLIWithResource("https://graph.windows.net")
	config.Resource = "https://graph.windows.net"
	authorizer, err := config.Authorizer()
	if err != nil {
		log.Println(err)
	}

	//authorizer, err := auth.NewAuthorizerFromCLI()
	//if err != nil {
	//	log.Println(err)
	//}

	grClient.Authorizer = authorizer

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	appsIterator, err := grClient.ListComplete(ctx, "")
	if err != nil{
		log.Println("Bad request when listing apps")
		log.Println(err)
	}

	for list := appsIterator; list.NotDone(); err = list.NextWithContext(ctx) {
		if err != nil {
			log.Fatalf("got error: %s", err)
		}

		printer.Data("Name: %s", *list.Value().DisplayName)
		for _, role := range *list.Value().AppRoles{
			printer.Data("Role: %s", role.Value)
		}
	}
}

