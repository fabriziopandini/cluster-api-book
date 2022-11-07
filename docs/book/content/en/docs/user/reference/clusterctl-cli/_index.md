---
title: "clustectl CLI"
weight: 50
description: |
  Documentation about the clusterctl CLI.
---

The `clusterctl` CLI tool handles the lifecycle of a Cluster API [management cluster].

The `clusterctl` command line interface is specifically designed for providing a simple "day 1 experience" and a
quick start with Cluster API. It automates fetching the YAML files defining [provider components] and installing them.

Additionally it encodes a set of best practices in managing providers, that helps the user in avoiding
mis-configurations or in managing day 2 operations such as upgrades.

## Avoiding GitHub rate limiting

While using providers hosted on GitHub, clusterctl is calling GitHub API which are rate limited; for normal usage free tier is enough but when using clusterctl extensively users might hit the rate limit.

To avoid rate limiting for the public repos set the `GITHUB_TOKEN` environment variable. To generate a token [follow this](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) documentation. The token only needs `repo` scope for clusterctl.

## Installing clusterctl
Instructions are available in the [Quick Start](../user/quick-start.md#install-clusterctl).

<!-- links -->
[management cluster]: ../reference/glossary.md#management-cluster
[provider components]: ../reference/glossary.md#provider-components

## Available commands