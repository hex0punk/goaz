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

type NSGState struct {
	AuditAll            bool
	Compact				bool
}

var (
	insecurePorts   = []string{ "22", "3389", "21", "20", "23"}
	nsgState      NSGState

	nsgComd = &cobra.Command{
		Use:   "nsg",
		Short: "nsg",
		Long:  `audit nsg`,
		Args: func(cmd *cobra.Command, args []string) error {
			if SubscriptionId == "" {
				return errors.New("please specify the subscription ID")
			}
			return nil
		},
		Run: func (cmd *cobra.Command, args []string) {
			if nsgState.AuditAll {
				nsgState.Audit()
			}
		},
	}
)

func init() {
	auditState = AuditSate{}
	nsgComd.Flags().BoolVarP(&nsgState.AuditAll, "Audit", "A", false, "-A")
	nsgComd.Flags().BoolVarP(&nsgState.Compact, "Compact", "C", false, "-C")

	rootCmd.AddCommand(nsgComd)
}


func (s *NSGState) Audit() {
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


	columns := []string{"ACCESS","DIRECTION","FROM ADDRESS", "FROM PORT", "TO ADDRESS", "TO PORT", "Insecure"}
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
				if rule.DestinationPortRange != nil{ // TODO: Check with a better properties
					ruleAssessment := getRuleAssessment(rule)
					if s.Compact && ruleAssessment == ""{
						continue
					}
					row := []string{string(rule.Access),string(rule.Direction), *rule.SourceAddressPrefix, *rule.SourcePortRange, *rule.DestinationAddressPrefix, *rule.DestinationPortRange, ruleAssessment}
					resultTable.Rows = append(resultTable.Rows, row)
				}
			}
			printer.PrintTable(false, &resultTable)
		}
	}
}

func getRuleAssessment(rule network.SecurityRule) string {
	if string(rule.Access) == "Deny"{
		return ""
	}
	if string(rule.Direction) == "Inbound" && *rule.SourceAddressPrefix == "*"{
		if *rule.DestinationPortRange == "*"{
			return "Insecure"
		}
		return "Public Access"
	}
	insecurePort := false
	for _, dp := range *rule.DestinationPortRanges{
		for _, p := range insecurePorts{
			if p == dp{
				insecurePort = true
				break
			}
		}
		if insecurePort{
			break
		}
	}
	if *rule.SourceAddressPrefix == "*" && insecurePort && string(rule.Direction) == "Inbound"{
		return "Insecure"
	}
	return ""
}
