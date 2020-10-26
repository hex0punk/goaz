package cmd

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hex0punk/goaz/utils"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
	"time"
)

type NSGState struct {
	AuditAll            bool
	Compact				bool
}

var (
	insecurePorts   = []int{ 22, 3389, 21, 20, 23}
	nsgState      NSGState

	nsgCmd = &cobra.Command{
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
	nsgCmd.Flags().BoolVarP(&nsgState.AuditAll, "Audit", "A", false, "-A")
	nsgCmd.Flags().BoolVarP(&nsgState.Compact, "Compact", "C", false, "-C")

	rootCmd.AddCommand(nsgCmd)
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


	columns := []string{"STATE","NAME","ACCESS","DIRECTION","FROM ADDRESS", "FROM PORT", "TO ADDRESS", "TO PORT", "EVAL	"}
	printer.Data("********Security Groups********\n")
	var sgList []network.SecurityGroup
	for _, rgName := range groups {
		sgIterator, err := sgClient.ListComplete(ctx, rgName)
		if err != nil{
			log.Println(err)
		}

		for list := sgIterator; list.NotDone(); err = list.NextWithContext(ctx) {
			if err != nil {
				log.Fatalf("got error: %s\n", err)
			}
			sgList = append(sgList, list.Value())
		}
	}
	for _, sg := range sgList{
		sgName := *sg.Name
		printer.InfoHeading("\t- Security Group Name: %s\n", sgName)
		rules := *sg.SecurityRules
		resultTable := printer.ResultTable{
			Columns: columns,
		}

		for _, rule := range rules{

			if rule.Name != nil{ // TODO: Check with a better properties
				ruleAssessment := getRuleAssessment(rule)
				if s.Compact && ruleAssessment == ""{
					continue
				}
				row := []string{string(strconv.FormatInt(int64(*rule.Priority),10)),*rule.Name,string(rule.Access),string(rule.Direction), *rule.SourceAddressPrefix, *rule.SourcePortRange, *rule.DestinationAddressPrefix, *rule.DestinationPortRange, ruleAssessment}
				resultTable.Rows = append(resultTable.Rows, row)
			}
		}
		if len(resultTable.Rows) > 0 {
			printer.PrintTable(&resultTable)
		}
	}
}

func getRuleAssessment(rule network.SecurityRule) string {
	if string(rule.Access) == "Deny"{
		return ""
	}

	insecurePort := false
	usesPortRange := false
	var firstLast []string

	if strings.Contains(*rule.DestinationPortRange, "-"){
		usesPortRange = true
		firstLast = strings.Split(*rule.DestinationPortRange, "-")
	}

	if usesPortRange && firstLast != nil{
		first, _ := strconv.Atoi(firstLast[0])
		last, _ := strconv.Atoi(firstLast[1])
		for i := first; i <= last; i++ {
			insecurePort = isPortInsecure(i)
			if insecurePort{
				break
			}
		}
	}

	if !usesPortRange{
		i, _ := strconv.Atoi(*rule.DestinationPortRange)
		insecurePort = isPortInsecure(i)
	}

	if string(rule.Direction) == "Inbound" && *rule.SourceAddressPrefix == "*"{
		if *rule.DestinationPortRange == "*" || insecurePort{
			//r,_ := rule.MarshalJSON()
			//log.Println(string(r))
			return "Insecure Port"
		}
		return "Public Access"
	}
	if *rule.SourceAddressPrefix == "*" && insecurePort && string(rule.Direction) == "Inbound"{
		return "Insecure"
	}
	return ""
}

func isPortInsecure(port int) bool{
	for _, p := range insecurePorts{
		if p == port{
			return true
		}
	}
	return false
}
