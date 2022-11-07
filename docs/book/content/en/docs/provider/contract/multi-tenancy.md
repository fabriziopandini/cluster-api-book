---
title: "Multi-tenancy"
weight: 40
description: |
  Documentation about provider's multy tenancy support.
---

Multi tenancy in Cluster API defines the capability of an infrastructure provider to manage different credentials, each
one of them corresponding to an infrastructure tenant.

In order to support multi tenancy, the following rule applies:

- Infrastructure providers MUST be able to manage different sets of credentials (if any)
- Providers SHOULD deploy and run any kind of webhook (validation, admission, conversion)
  following Cluster API codebase best practices for the same release.
- Providers MUST create and publish a `{type}-component.yaml` accordingly.
