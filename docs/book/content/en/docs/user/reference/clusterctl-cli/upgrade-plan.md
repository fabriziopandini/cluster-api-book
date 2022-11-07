---
title: "clusterctl upgrade plan"
linkTitle: "upgrade plan"
description: |
  Provide a list of recommended target versions for upgrading Cluster API providers in a management cluster.
---
# upgrade plan

The `clusterctl upgrade plan` command can be used to identify possible targets for upgrades.


```bash
clusterctl upgrade plan
```

Produces an output similar to this:

```bash
Checking cert-manager version...
Cert-Manager will be upgraded from "v1.5.0" to "v1.5.3"

Checking new release availability...

Management group: capi-system/cluster-api, latest release available for the v1beta1 API Version of Cluster API (contract):

NAME                    NAMESPACE                           TYPE                     CURRENT VERSION   NEXT VERSION
bootstrap-kubeadm       capi-kubeadm-bootstrap-system       BootstrapProvider        v0.4.0           v1.0.0
control-plane-kubeadm   capi-kubeadm-control-plane-system   ControlPlaneProvider     v0.4.0           v1.0.0
cluster-api             capi-system                         CoreProvider             v0.4.0           v1.0.0
infrastructure-docker   capd-system                         InfrastructureProvider   v0.4.0           v1.0.0
```
You can now apply the upgrade by executing the following command:
```bash
   clusterctl upgrade apply --contract v1beta1
```

The output contains the latest release available for each API Version of Cluster API (contract)
available at the moment.

{{< alert >}}

`clusterctl upgrade plan` does not display pre-release versions by default. For
example, if a provider has releases `v0.7.0-alpha.0` and `v0.6.6` available, the latest
release available for upgrade will be `v0.6.6`.

{{< /alert >}}