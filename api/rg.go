package api

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/go-autorest/autorest"
	"log"
	"time"
)

func GetResourceGroups(authorizer autorest.Authorizer, subId string) []string {
	groupsClient := resources.NewGroupsClient(subId)
	groupsClient.Authorizer = authorizer

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	groupIterator, err := groupsClient.ListComplete(ctx, "", nil)
	var groups []string
	for list := groupIterator; list.NotDone(); err = list.NextWithContext(ctx) {
		if err != nil {
			log.Fatalf("got error: %s\n", err)
		}
		rgName := *list.Value().Name
		groups = append(groups, rgName)
	}

	return groups
}