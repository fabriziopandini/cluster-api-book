# Cluster API book

In this project, I'm playing around [Hugo][] and [Docsy][] for usage on Cluster API book.

Preview visible on [Netlify][] at [](https://cluster-api-book-test.netlify.app/)

[Docsy]: https://github.com/google/docsy
[Hugo]: https://gohugo.io/
[Netlify]: https://netlify.com

## Testing locally

Serving the site locally:

```bash
make serve-book
```

Run the lint checker:

```bash
make verify-markdown-link
```

## Where we are

| Surce                                                                        | Target                                                                                                     | Notes                          |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ------------------------------ |
| Home (new page)                                                              | \_index.md                                                                                                 | TODO, providers, user's quotes |
| About (new page)                                                             | \about\_index.md                                                                                           | TODO                           |
| Community (new page)                                                         | \community\_index.md                                                                                       | TODO                           |
| \CONTRIBUTING.md                                                             | \docs\develop\contributing.md                                                                              | TODO, include                  |
| \REVIEWING.md                                                                | \docs\develop\review\reviewing.md                                                                          | TODO, include                  |
| \clusterctl\commands\additional-commands.md                                  | \docs\user\reference\clusterctl-cli\backup.md & config-repositories.md & help.md & restore.md & version.md | OK 10-21-2022                  |
| \clusterctl\commands\alpha-rollout.md                                        | \docs\user\reference\clusterctl-cli\alpha-rollout.md                                                       | OK 10-21-2022                  |
| \clusterctl\commands\alpha-topology-plan.md                                  | \docs\user\reference\clusterctl-cli\alpha-topology-plan.md                                                 | OK 10-21-2022                  |
| \clusterctl\commands\commands.md                                             | \docs\user\reference\clusterctl-cli\_index.md                                                              | OK 10-21-2022                  |
| \clusterctl\commands\completion.md                                           | \docs\user\reference\clusterctl-cli\completion.md                                                          | OK 10-21-2022                  |
| \clusterctl\commands\delete.md                                               | \docs\user\reference\clusterctl-cli\delete.md                                                              | OK 10-21-2022                  |
| \clusterctl\commands\describe-cluster.md                                     | \docs\user\reference\clusterctl-cli\describe-cluster.md                                                    | OK 10-21-2022                  |
| \clusterctl\commands\generate-cluster.md                                     | \docs\user\reference\clusterctl-cli\generate-cluster.md                                                    | OK 10-21-2022                  |
| \clusterctl\commands\generate-provider.md                                    | \docs\user\reference\clusterctl-cli\generate-provider.md                                                   | OK 10-21-2022                  |
| \clusterctl\commands\generate-yaml.md                                        | \docs\user\reference\clusterctl-cli\generate-yaml.md                                                       | OK 10-21-2022                  |
| \clusterctl\commands\get-kubeconfig.md                                       | \docs\user\reference\clusterctl-cli\get-kubeconfig.md                                                      | OK 10-21-2022                  |
| \clusterctl\commands\init.md                                                 | \docs\user\reference\clusterctl-cli\init.md                                                                | OK 10-21-2022                  |
| \clusterctl\commands\move.md                                                 | \docs\user\reference\clusterctl-cli\move.md                                                                | OK 10-21-2022                  |
| \clusterctl\commands\upgrade.md                                              | \docs\user\reference\clusterctl-cli\upgrade-apply.md & upgrade-plan.md                                     | OK 10-21-2022                  |
| \clusterctl\configuration.md                                                 | \docs\user\reference\clusterctl-configuration.md                                                           | OK 10-21-2022                  |
| \clusterctl\developers.md                                                    | \docs\develop\test\clusterctl.md                                                                           | OK 10-30-2022                  |
| \clusterctl\overview.md                                                      | \docs\user\reference\clusterctl-cli\_index.md                                                              | OK 10-21-2022                  |
| \clusterctl\provider-contract.md                                             | \docs\provider\contract\clusterctl.md                                                                      | OK 10-30-2022                  |
| \code-of-conduct.md                                                          | \_index.md & \about\_index.md & \community\_index.md & \docs\develop\_index.md                             |                                |
| \developer\architecture\controllers\bootstrap.md                             | \docs\provider\contract\bootstrap-configuration.md                                                         |                                |
| \developer\architecture\controllers\cluster-resource-set.md                  | \docs\develop\controllers\cluster-resource-set.md                                                          |                                |
| \developer\architecture\controllers\cluster-topology.md                      | \docs\develop\controllers\cluster-topology.md                                                              |                                |
| \developer\architecture\controllers\cluster.md                               | \docs\develop\controllers\cluster.md                                                                       |                                |
| \developer\architecture\controllers\control-plane.md                         | \docs\provider\contract\control-plane.md                                                                   |                                |
| \developer\architecture\controllers\machine-deployment.md                    | \docs\develop\controllers\machine-deployment.md                                                            |                                |
| \developer\architecture\controllers\machine-health-check.md                  | \docs\develop\controllers\machine-health-check.md                                                          |                                |
| \developer\architecture\controllers\machine-pool.md                          | \docs\develop\controllers\machine-pool.md                                                                  |                                |
| \developer\architecture\controllers\machine-set.md                           | \docs\develop\controllers\machine-set.md                                                                   |                                |
| \developer\architecture\controllers\machine.md                               | \docs\develop\controllers\machine.md                                                                       |                                |
| \developer\architecture\controllers\metadata-propagation.md                  | \docs\user\reference\metadata-propagation.md                                                               |                                |
| \developer\architecture\controllers\multi-tenancy.md                         | \docs\provider\contract\multi-tenancy.md                                                                   | OK 10-30-2022                  |
| \developer\architecture\controllers\support-multiple-instances.md            | \docs\provider\contract\multiple-instances.md                                                              | OK 10-30-2022                  |
| \developer\architecture\controllers.md                                       | \docs\develop\controllers\_index.md                                                                        |                                |
| \developer\crd-relationships.md                                              | \docs\user\reference\crds.md                                                                               |                                |
| \developer\e2e.md                                                            | \docs\develop\test\e2e-tests.md                                                                            | OK 10-30-2022                  |
| \developer\guide.md                                                          | \docs\user\_index.md & prerequistes.md                                                                     |                                |
| \developer\logging.md                                                        | \docs\develop\controllers\logging.md                                                                       | OK 10-21-2022                  |
| \developer\providers\bootstrap.md                                            | \docs\provider\contract\bootstrap-configuration.md                                                         |                                |
| \developer\providers\cluster-infrastructure.md                               | \docs\provider\contract\cluster-infrastructure.md                                                          |                                |
| machine-pool-infrastructure (new page)                                       | \docs\provider\contract\machine-pool-infrastructure.md                                                     |                                |
| \developer\providers\contracts.md                                            | \docs\provider\contract\crd.md                                                                             |                                |
| \developer\providers\implementers-guide\building_running_and_testing.md      | \docs\provider\implementers-guide\build-run-test.md                                                        |                                |
| \developer\providers\implementers-guide\configure.md                         | \docs\provider\implementers-guide\config.md                                                                |                                |
| \developer\providers\implementers-guide\controllers_and_reconciliation.md    | \docs\provider\implementers-guide\controllers.md                                                           |                                |
| \developer\providers\implementers-guide\create_api.md                        | \docs\provider\implementers-guide\create-api.md                                                            |                                |
| \developer\providers\implementers-guide\generate_crds.md                     | \docs\provider\implementers-guide\create-repo-and-crd.md                                                   |                                |
| \developer\providers\implementers-guide\naming.md                            | \docs\provider\implementers-guide\naming.md                                                                |                                |
| \developer\providers\implementers-guide\overview.md                          | \developer\providers\implementers-guide\_index.md                                                          |                                |
| \developer\providers\implementers.md                                         | \developer\providers\implementers-guide\_index.md                                                          |                                |
| \developer\providers\machine-infrastructure.md                               | \docs\provider\contract\machine-infrastructure.md                                                          |                                |
| \developer\providers\v0.3-to-v0.4.md                                         | \docs\provider\version-migration\v0.3-to-v0.4.md                                                           |                                |
| \developer\providers\v0.4-to-v1.0.md                                         | \docs\provider\version-migration\v0.4-to-v1.0.md                                                           |                                |
| \developer\providers\v1.0-to-v1.1.md                                         | \docs\provider\version-migration\v1.0-to-v1.1.md                                                           |                                |
| \developer\providers\v1.1-to-v1.2.md                                         | \docs\provider\version-migration\v1.1-to-v1.2.md                                                           |                                |
| \developer\providers\v1.2-to-v1.3.md                                         | \docs\provider\version-migration\v1.2-to-v1.3.md                                                           |                                |
| \developer\providers\version-migration.md                                    | \docs\provider\version-migration\_index.md                                                                 |                                |
| \developer\providers\webhooks.md                                             | \docs\provider\implementers-guide\webhooks.md                                                              |                                |
| \developer\repository-layout.md                                              | \docs\develop\repository-layout.md                                                                         |                                |
| \developer\testing.md                                                        | \docs\develop\test\testing.md                                                                              | OK 10-30-2022                  |
| \developer\tilt.md                                                           | \docs\develop\dev\tilt.md                                                                                  | OK 10-30-2022                  |
| \introduction.md                                                             | \_index.md                                                                                                 |                                |
| \reference\api_reference.md                                                  | \docs\user\reference\crds.md & golang.md                                                                   | OK 10-21-2022                  |
| \reference\glossary.md                                                       | \docs\user\reference\glossary.md                                                                           | OK 10-30-2022                  |
| \reference\jobs.md                                                           | \docs\develop\test\jobs.md                                                                                 | OK 10-30-2022                  |
| \reference\ports.md                                                          | \docs\provider\contract\manager.md                                                                         |                                |
| \reference\providers.md                                                      | \_index.md                                                                                                 |                                |
| \reference\reference.md                                                      | \docs\user\reference\_index.md                                                                             |                                |
| \reference\versions.md                                                       | \docs\user\reference\version.md                                                                            |                                |
| labels-annotations (new page/PR in flight)                                   | \docs\user\reference\labels-annotations.md                                                                 |                                |
| \roadmap.md                                                                  |                                                                                                            | TBD                            |
| \security\index.md                                                           |                                                                                                            |                                |
| \security\pod-security-standards.md                                          |                                                                                                            |                                |
| \tasks\automated-machine-management\autoscaling.md                           |                                                                                                            |                                |
| \tasks\automated-machine-management\healthchecking.md                        |                                                                                                            |                                |
| \tasks\automated-machine-management\index.md                                 |                                                                                                            |                                |
| \tasks\automated-machine-management\scaling.md                               |                                                                                                            |                                |
| \tasks\bootstrap\index.md                                                    |                                                                                                            |                                |
| \tasks\bootstrap\kubeadm-bootstrap.md                                        |                                                                                                            |                                |
| \tasks\bootstrap\microk8s-bootstrap.md                                       |                                                                                                            |                                |
| \tasks\certs\generate-kubeconfig.md                                          |                                                                                                            |                                |
| \tasks\certs\index.md                                                        |                                                                                                            |                                |
| \tasks\certs\using-custom-certificates.md                                    |                                                                                                            |                                |
| \tasks\control-plane\index.md                                                |                                                                                                            |                                |
| \tasks\control-plane\kubeadm-control-plane.md                                |                                                                                                            |                                |
| \tasks\control-plane\microk8s-control-plane.md                               |                                                                                                            |                                |
| \tasks\experimental-features\tasks\change-clusterclass.md                    |                                                                                                            |                                |
| \tasks\experimental-features\tasks\index.md                                  |                                                                                                            |                                |
| \tasks\experimental-features\tasks\operate-cluster.md                        |                                                                                                            |                                |
| \tasks\experimental-features\tasks\write-clusterclass.md                     |                                                                                                            |                                |
| \tasks\experimental-features\cluster-resource-set.md                         |                                                                                                            |                                |
| \tasks\experimental-features\experimental-features.md                        |                                                                                                            |                                |
| \tasks\experimental-features\ignition.md                                     |                                                                                                            |                                |
| \tasks\experimental-features\machine-pools.md                                |                                                                                                            |                                |
| \tasks\experimental-features\runtime-sdk\deploy-runtime-extension.md         |                                                                                                            |                                |
| \tasks\experimental-features\runtime-sdk\implement-extensions.md             |                                                                                                            |                                |
| \tasks\experimental-features\runtime-sdk\implement-lifecycle-hooks.md        |                                                                                                            |                                |
| \tasks\experimental-features\runtime-sdk\implement-topology-mutation-hook.md |                                                                                                            |                                |
| \tasks\experimental-features\runtime-sdk\index.md                            |                                                                                                            |                                |
| \tasks\external-etcd.md                                                      |                                                                                                            |                                |
| \tasks\index.md                                                              |                                                                                                            |                                |
| \tasks\updating-machine-templates.md                                         |                                                                                                            |                                |
| \tasks\upgrading-cluster-api-versions.md                                     |                                                                                                            |                                |
| \tasks\upgrading-clusters.md                                                 |                                                                                                            |                                |
| \tasks\using-kustomize.md                                                    |                                                                                                            |                                |
| \concepts.md                                                                 | \docs\user\concepts.md                                                                                     |                                |
| \personas.md                                                                 | \docs\user\personas.md                                                                                     |                                |
| \quick-start.md                                                              | \try\_index.md                                                                                             |                                |
| \troubleshooting.md                                                          | \docs\user\troubleshooting.md                                                                              |                                |
| search (new page)                                                            | \search.md                                                                                                 |                                |

## TODO

There is still a long way to do...

### Look and feel

- [ ] on small devices convert navbar menu to drop list
- [ ] prevent .td-sidebar to go under the logo when scrolling down
- [ ] 404
- [ ] give credits to https://github.com/rvanhorn/hugo-dynamic-tabs

### Missing features

- [ ] Release link
- [ ] Embed github (iframe??)

### Content

- [ ] complete work on "home", "about" and "community"
  - [ ] think to key messages in home (multi cloud, like Kubernetes, manage etc, might be an example)
  - [ ] ensure "about" and "community" can have look and feel like home page at top and content below


