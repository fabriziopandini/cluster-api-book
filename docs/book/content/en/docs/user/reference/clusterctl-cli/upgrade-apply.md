---
title: "clusterctl upgrade apply"
linkTitle: "upgrade apply"
description: |
  Apply new versions of Cluster API core and providers in a management cluster.
---

After choosing the desired option for the upgrade, you can run the following
command to upgrade all the providers in the management cluster. This upgrades
all the providers to the latest stable releases.

```bash
clusterctl upgrade apply --contract v1beta1
```

The upgrade process is composed by three steps:

* Check the cert-manager version, and if necessary, upgrade it.
* Delete the current version of the provider components, while preserving the namespace where the provider components
  are hosted and the provider's CRDs.
* Install the new version of the provider components.

Please note that clusterctl does not upgrade Cluster API objects (Clusters, MachineDeployments, Machine etc.); upgrading
such objects are the responsibility of the provider's controllers.

{{< alert type="warning" >}}

The current implementation of the upgrade process does not preserve controllers flags that are not set through the
components YAML/at the installation time.

User is required to re-apply flag values after the upgrade completes.

{{< /alert >}}

In order to upgrade to a provider's pre-release version, use
the following syntax:

```bash
clusterctl upgrade apply \
    --core capi-system/cluster-api:v1.0.0 \
    --bootstrap capi-kubeadm-bootstrap-system/kubeadm:v1.0.0 \
    --control-plane capi-kubeadm-control-plane-system/kubeadm:v1.0.0 \
    --infrastructure capd-system/docker:v1.0.0-rc.0
```

In this case, all the provider's versions must be explicitly stated.

