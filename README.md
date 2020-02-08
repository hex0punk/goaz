# Auditing Azure with goaz

Goaz is a simple application meant to help researchers and blue teams audit azure. At the time it only works with Azure storage.

## Getting Started

Currently the only way to authenticate is to log in to the Azure CLI and run `goaz`. `goaz` will then use the current CLI authentication values to do its job.

All commands require that you enter a `subscriptionId` value so that goaz knows which subscription to work with.

### Storage

#### Auditing

Goaz checks the following types of Azure storage:

- Blobs
- File Shares
- Storage Queues

To perform an audit of all storage types listed above type the following:

```shell
goaz audit --subscriptionId <subscription ID> -A
```

#### Stalking Queues

Goaz can also monitor storage queues by "peeking" into any given queue. Note that this does not remove messages from the queue. Use this functionality sparingly, as peeking into a queue can result in additional charges on your Azure account.

To stalk a message queue type the following:

```shell
go run main.go stalk -q --subscriptionId <subscription ID> --account <storage account name> -name outqueue --key <storage account key>
```

### Network Security Groups

Goaz checks for insecure security group settings:

```shell
goaz nsg --subscriptionId <subscription ID> -A
```

