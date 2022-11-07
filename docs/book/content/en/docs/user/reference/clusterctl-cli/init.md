---
title: "clustectl init"
linkTitle: "init"
description: |
  Initialize a management cluster.
---

The `clusterctl init` command installs the Cluster API components and transforms the Kubernetes cluster
into a management cluster.

This document provides more detail on how `clusterctl init` works and on the supported options for customizing your
management cluster.

## Defining the management cluster

The `clusterctl init` command accepts in input a list of providers to install.

You can use the `clusterctl config repositories` command to get a list of supported providers and their repository configuration.

If the provider of your choice is missing, you can customize the list of supported providers by using the
[clusterctl configuration](../configuration.md) file.

### Automatically installed providers

The `clusterctl init` command automatically adds the `cluster-api` core provider, the `kubeadm` bootstrap provider, and
the `kubeadm` control-plane provider to the list of providers to install. 

This allows users to use a concise command syntax for initializing a management cluster.
For example, to get a fully operational management cluster with the `aws` infrastructure provider, the `cluster-api` core provider, the `kubeadm` bootstrap, and the `kubeadm` control-plane provider, use the command:

`clusterctl init --infrastructure aws`

{{< alert >}}

The `cluster-api` core provider, the `kubeadm` bootstrap provider, and the `kubeadm` control-plane provider are automatically installed only if:
- The user doesn't explicitly require to install a core/bootstrap/control-plane provider using the `--core` flag, the `--bootstrap` flag or the `--control-plane` flags;
- There is not an instance of a CoreProvider already installed in the cluster;

Please note that the second rule allows to execute `clusterctl init` more times: the first call actually initializes
the management cluster, while the subsequent calls can be used to add more providers.

{{< /alert >}}

{{< alert type="warning" >}}

The `clusterctl init` command forbids users from installing two instances of the *same* provider.

{{< /alert >}}

To skip automatic provider installation use  `--bootstrap "-"` or  `--control-plane "-"`;
it is not possible to skip automatic installation of the `cluster-api` core provider.

### Provider version

The `clusterctl init` command by default installs the latest version available
for each selected provider; by default pre-release versions are not installed. For
example, if a provider has releases `v0.7.0-alpha.0` and `v0.6.6`, the latest
release installed will be `v0.6.6`.

If required, it possible to install a specific version of a provider by appending a version tag to the provider name, e.g. `aws:v0.4.1`; this approach can be used to install pre-release versions, e.g. `vsphere:v0.7.0-alpha.0`. 

{{< alert >}}

Pinning the version provides better control over what clusterctl chooses to install
(usually required in an enterprise environment). Version pinning should always be used when using [image overrides](../configuration.md#image-overrides), or when relying on internal repositories with a separated
software supply chain, or a custom versioning schema.

{{< /alert >}}

#### Target namespace

The `clusterctl init` command by default installs each provider in the default target namespace defined by each provider, e.g. `capi-system` for the Cluster API core provider.

See the provider documentation for more details.

If required, it possible to change the target namespace by using the `--target-namespace` flag.

Please, note that the `--target-namespace` flag applies to all the providers to be installed during a `clusterctl init` operation.

## Provider repositories

To access provider specific information, such as the components YAML to be used for installing a provider,
`clusterctl init` accesses the **provider repositories**, that are well-known places where the release assets for
a provider are published.

If, for any reasons, the user wants to replace the assets available on a provider repository with a locally available asset,
the user is required to save the file under `$HOME/.cluster-api/overrides/<provider-label>/<version>/<file-name.yaml>`.

```bash
$HOME/.cluster-api/overrides/infrastructure-aws/v0.5.2/infrastructure-components.yaml
```

See [clusterctl configuration](../configuration.md) for more info about provider repository configurations and overrides.

## Variable substitution
Providers can use variables in the components YAML published in the provider's repository.

Users can refer to the provider documentation for the list of variables to be set or use the
`clusterctl generate provider --<provider-type> <provider-name> --describe` command to get a list of expected variable names.

During `clusterctl init`, those variables are replaced with environment variables or with variables read from the
[clusterctl configuration](../configuration.md).

{{< alert type="warning" >}}

The user should ensure the variables required by a provider are set in advance.

{{< /alert >}}

## Additional information

When installing a provider, the `clusterctl init` command executes a set of steps to simplify
the lifecycle management of the provider's components.

* All the provider's components are labeled, so they can be easily identified in
subsequent moments of the provider's lifecycle, e.g. upgrades.

 ```bash
 labels:
 - clusterctl.cluster.x-k8s.io: ""
 - cluster.x-k8s.io/provider: "<provider-name>"
 ```

* An additional `Provider` object is created in the target namespace where the provider is installed.
This object keeps track of the provider version, and other useful information
for the inventory of the providers currently installed in the management cluster.

<aside class="note warning">

<h1>Warning</h1>

The `clusterctl.cluster.x-k8s.io` labels, the `cluster.x-k8s.io/provider` labels and the `Provider` objects MUST NOT be altered.
If this happens, there are no guarantees about the proper functioning of `clusterctl`.

</aside>

## Cert-manager

Cluster API providers require a cert-manager version supporting the `cert-manager.io/v1` API to be installed in the cluster.

While doing init, clusterctl checks if there is a version of cert-manager already installed. If not, clusterctl will
install a default version (currently cert-manager v1.10.0). See [clusterctl configuration](../configuration.md) for
available options to customize this operation.

<aside class="note warning">

<h1>Warning</h1>

Please note that, if clusterctl installs cert-manager, it will take care of its lifecycle, eventually upgrading it
during clusterctl upgrade. Instead, if cert-manager is provided by the users, the user is responsible for
upgrading this component when required.

</aside>

## Avoiding GitHub rate limiting

Follow [this](../overview.md#avoiding-github-rate-limiting)