---
title: "Testing clusterctl with a local repository"
weight: 20
description: |
  Documentation about how to test clusterctl using a local repository.
---

When developing clusterctl it is a common practice to use local repository whenever
it is required to test deploying artifacts/versions of the providers not yet published
in the official repositories. 

## Prerequisites

* A Cluster API development setup (go, git, kind v0.9 or newer, Docker v19.03 or newer etc.)
* A local clone of the Cluster API GitHub repository
* A local clone of the GitHub repositories for the providers you want to install

## Build clusterctl

From the root of the local copy of Cluster API, you can build the `clusterctl` binary by running:

```bash
make clusterctl
```

The output of the build is saved in the `bin/` folder; In order to use it you have to specify
the full path, create an alias or copy it into a folder under your `$PATH`.

## Use local artifacts

Clusterctl by default uses artifacts published in the [providers repositories];
during the development workflow you may want to use artifacts from your local workstation.

There are two options to do so:

* Use the [overrides layer], when you want to override a single published artifact with a local one.
* Create a local repository, when you want to avoid using published artifacts and use the local ones instead.

If you want to create a local artifact, follow these instructions:

### Build artifacts locally

In order to build artifacts for the CAPI core provider, the kubeadm bootstrap provider, the kubeadm control plane provider and the docker infrastructure provider:

```bash
make docker-build REGISTRY=gcr.io/k8s-staging-cluster-api PULL_POLICY=IfNotPresent
```

### Create a clusterctl-settings.json file

Next, create a `clusterctl-settings.json` file and place it in your local copy
of Cluster API. This file will be used by [create-local-repository.py](#create-the-local-repository). Here is an example:

```yaml
{
  "providers": ["cluster-api","bootstrap-kubeadm","control-plane-kubeadm", "infrastructure-aws", "infrastructure-docker"],
  "provider_repos": ["../cluster-api-provider-aws"]
}
```

**providers** (Array[]String, default=[]): A list of the providers to enable.
See [available providers](#available-providers) for more details.

**provider_repos** (Array[]String, default=[]): A list of paths to all the providers you want to use. Each provider must have
a `clusterctl-settings.json` file describing how to build the provider assets.

### Create the local repository

Run the create-local-repository hack from the root of the local copy of Cluster API:

```bash
cmd/clusterctl/hack/create-local-repository.py
```

The script reads from the source folders for the providers you want to install, builds the providers' assets,
and places them in a local repository folder located under `$HOME/.cluster-api/dev-repository/`.
Additionally, the command output provides you the `clusterctl init` command with all the necessary flags.
The output should be similar to:

```bash
clusterctl local overrides generated from local repositories for the cluster-api, bootstrap-kubeadm, control-plane-kubeadm, infrastructure-docker, infrastructure-aws providers.
in order to use them, please run:

clusterctl init \
   --core cluster-api:v0.3.8 \
   --bootstrap kubeadm:v0.3.8 \
   --control-plane kubeadm:v0.3.8 \
   --infrastructure aws:v0.5.0 \
   --infrastructure docker:v0.3.8 \
   --config ~/.cluster-api/dev-repository/config.yaml
```

As you might notice, the command is using the `$HOME/.cluster-api/dev-repository/config.yaml` config file,
containing all the required setting to make clusterctl use the local repository.


