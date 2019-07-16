package storageapi

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/azure-storage-file-go/azfile"
	"github.com/Azure/azure-storage-queue-go/azqueue"
	"net/url"
)

var (
	blobFormatString = `https://%s.blob.core.windows.net`
	fileFormatString = `https://%s.file.core.windows.net`
	queueFormatString = `https://%s.queue.core.windows.net`
	messageFormatString = `https://%s.queue.core.windows.net/%s/messages`
)

func GetBlobContainerURL(accountName string, containerName string, key string) azblob.ContainerURL {
	c, _ := azblob.NewSharedKeyCredential(accountName, key)
	p := azblob.NewPipeline(c, azblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf(blobFormatString, accountName))
	service := azblob.NewServiceURL(*u, p)
	container := service.NewContainerURL(containerName)
	return container
}

func GetFilesURL(accountName string, directoryName string, key string) azfile.DirectoryURL {
	c, _ := azfile.NewSharedKeyCredential(accountName, key)
	p := azfile.NewPipeline(c, azfile.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf(fileFormatString, accountName))
	service := azfile.NewShareURL(*u, p)
	directory := service.NewDirectoryURL(directoryName)
	return directory
}

func PeekMessages(ctx context.Context, accountName string, queueName string, key string, max int32) *azqueue.PeekedMessagesResponse {
	c, _ := azqueue.NewSharedKeyCredential(accountName, key)
	p := azqueue.NewPipeline(c, azqueue.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf(messageFormatString, accountName, queueName))
	messagesUrl := azqueue.NewMessagesURL(*u, p)
	messages, err := messagesUrl.Peek(ctx, max)
	if err != nil{
		fmt.Println(err.Error())
	}
	return messages
}

func ListQueues(ctx context.Context, accountName string, key string) *azqueue.ListQueuesSegmentResponse{
	c, _ := azqueue.NewSharedKeyCredential(accountName, key)
	p := azqueue.NewPipeline(c, azqueue.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf(queueFormatString, accountName))
	service := azqueue.NewServiceURL(*u, p)
	queues, _ := service.ListQueuesSegment(ctx, azqueue.Marker{}, azqueue.ListQueuesSegmentOptions{})
	return queues
}

func GetShareDirectories(ctx context.Context, accountName string, key string) *azfile.ListSharesResponse{
	c, _ := azfile.NewSharedKeyCredential(accountName, key)
	p := azfile.NewPipeline(c, azfile.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf(fileFormatString, accountName))
	service := azfile.NewServiceURL(*u, p)
	directories, _ := service.ListSharesSegment(ctx, azfile.Marker{}, azfile.ListSharesOptions{})
	return directories
}

func ListBlobs(ctx context.Context, accountName, containerName string, key string) (*azblob.ListBlobsFlatSegmentResponse, error) {
	c := GetBlobContainerURL(accountName, containerName, key)
	return c.ListBlobsFlatSegment(
		ctx,
		azblob.Marker{},
		azblob.ListBlobsSegmentOptions{
			Details: azblob.BlobListingDetails{
				Snapshots: true,
			},
		})
}

func ListFiles(ctx context.Context, accountName, directoryName string, key string) (*azfile.ListFilesAndDirectoriesSegmentResponse, error) {
	c := GetFilesURL(accountName, directoryName, key)
	return c.ListFilesAndDirectoriesSegment(
		ctx,
		azfile.Marker{},
		azfile.ListFilesAndDirectoriesOptions{})
}

func BlobURL(account string, container string, name string) string{
	return fmt.Sprintf(blobFormatString, account) + fmt.Sprintf("/%s/%s", container, name)
}

func FileURL(account string, share string, name string) string{
	return fmt.Sprintf(fileFormatString, account) + fmt.Sprintf("/%s/%s", share, name)
}