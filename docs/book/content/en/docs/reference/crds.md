---
title: "API types"
weight: 20
description: |
  Documentation about API types from Cluster API Custom Resource Definitions.
migration: |
  Last import date: 10-21-2022 
  NOTE: golang reference from <old-book>/reference/api_reference has been moved into golang.md
aliases:
- /reference/api_reference
- /reference/api_reference.html
- /developer/crd-relationships
- /developer/crd-relationships.html
---

## API types definitions

API types from [Custom Resource Definitions](https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api) (courtesy of https://doc.crds.dev) .

## API types relationships

There are many resources that appear in the Cluster API. In this section, we use diagrams to illustrate the most common relationships between Cluster API resources.

{{< alert >}}

The straight lines represent "management". For example, "MachineSet manages Machines". The dotted line represents "reference". For example, "Machine's `spec.infrastructureRef` field references FooMachine".

The direction of the arrows indicates the direction of "management" or "reference". For example, "the relationship between MachineSet and Machine is management from MachineSet to Machine", so the arrow points from MachineSet to Machine.

{{< /alert >}}

### KubeadmControlPlane machines

![](/images/kubeadm-control-plane-machines-resources.png)

### Worker machines

![](/images/worker-machines-resources.png)
