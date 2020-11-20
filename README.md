# Auditing Azure with goaz

goaz is a simple application meant to help researchers and blue teams audit azure. The following resources are supported:

- Azure Kubernetes Services
- Virtual Machine Scale Sets
- Azure storage (containers, blobs, file shares, queues)
- Service Bus
- Key Vault
- Azure Public Addresses
- Network Security Groups

## Getting Started

Currently, the only way to authenticate is to log in to the Azure CLI using `az login` and run `goaz`. `goaz` will then use the current CLI authentication values to do its job.

All commands require that you enter a `subscriptionId` value so that goaz knows which subscription to work with.

## Supported checks

### Storage

Goaz checks the following types of Azure storage and verifies that secure transfers are enabled and that Firewall and VNET restrictions are in place. It also flags any storage resource with a public access type other than none.

- Blobs
- File Shares
- Storage Queues

To perform an audit of all storage types listed above type the following:

```shell
goaz storage --subscriptionId <subscription ID> -A
```

You can also specify the resource group if desired:

```shell
goaz storage --subscriptionId <subscription ID> --resourceGroup <resource group name> -A
```

#### Stalking Queues

Goaz can also monitor storage queues by "peeking" into any given queue. Note that this does not remove messages from the queue. Use this functionality sparingly, as peeking into a queue can result in additional charges on your Azure account.

To stalk a message queue type the following:

```shell
goaz stalk -q --subscriptionId <subscription ID> --account <storage account name> -name outqueue --key <storage account key>
```

### Virtual Machine Scale sets

Goaz will look for issues due to missing Azure Disk Encryption (ADE), and will verify that boot diagnostics are turned on. It will also flag VMSS that are not configured with security groups.

```shell
goaz vms --subscriptionId <subscription ID>
```

### Azure Kubernetes Services

At the moment, goaz will only list basic information for AKS, including the URL for the k8s API

```shell
goaz aks --subscriptionId <subscription ID>
```

### Message Bus

Goaz checks whether redundancy is enabled and whether VNET and Firewall rules are in place restricting public access to the queues.

```shell
goaz sbus --subscriptionId <subscription ID>
```

### Azure Key Vault

Goaz checks that Key Vaults are configured with Firewall rules and their access restricted by VNETs. It will also detect whether keys are used for deployments or disk encryption.

```shell
goaz kv --subscriptionId <subscription ID>
```

### Network

Provided by `goaz net`

#### Network Security Groups

Goaz checks for insecure security group settings:

```shell
goaz net nsg --subscriptionId <subscription ID>
```

#### Public IPs

Goaz checks for Azure public IPs and verifies DDoS protections are in place:

```shell
goaz net pips --subscriptionId <subscription ID>
```
