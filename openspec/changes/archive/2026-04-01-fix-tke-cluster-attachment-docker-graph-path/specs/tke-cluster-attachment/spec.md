## MODIFIED Requirements

### Requirement: docker_graph_path uses platform default when not explicitly set

The `docker_graph_path` field in both `worker_config` and `worker_config_overrides` blocks SHALL use `Computed: true` instead of a hardcoded `Default` value.

When a user does not set `docker_graph_path`, the platform-determined value SHALL be read back from the API and stored in state.

Existing resources that already have `/var/lib/docker` stored in their Terraform state SHALL NOT generate a plan diff after this change.

#### Scenario: New resource without docker_graph_path set

- **WHEN** a user creates a `tencentcloud_kubernetes_cluster_attachment` resource without specifying `docker_graph_path`
- **THEN** Terraform SHALL apply the resource and store the platform-returned value (e.g. `/var/lib/containerd`) in state

#### Scenario: Existing resource with /var/lib/docker in state

- **WHEN** a user runs `terraform plan` on an existing resource that has `/var/lib/docker` recorded in state
- **THEN** Terraform SHALL NOT show a diff for `docker_graph_path` and SHALL NOT propose recreation of the resource
