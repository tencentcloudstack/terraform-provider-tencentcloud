Use this data source to query detailed information of TKE addon values.

-> **Note:** `values` is the rendering result of the parameters you configured merged with the
default parameters of the chart. It is **not** empty for an addon that is not installed in the
cluster, as long as the chart exists in the addon repository: in that case `values` is exactly
equal to `default_values` and contains no cluster information (for example `global.cluster.id`
is an empty string). To check whether an addon is actually installed, use the
`tencentcloud_kubernetes_addons` data source instead.

Example Usage

```hcl
data "tencentcloud_tke_addon_values" "example" {
  cluster_id = "cls-5yezvaxo"
  addon_name = "ip-masq-agent"
}

locals {
  current_values   = jsondecode(data.tencentcloud_tke_addon_values.example.values)
  available_values = jsondecode(data.tencentcloud_tke_addon_values.example.default_values)
}

output "addon_current_values" {
  value = local.current_values
}

output "addon_available_values" {
  value = local.available_values
}
```

Example Usage of tuning an installed addon

```hcl
data "tencentcloud_tke_addon_values" "rdma_agent" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  addon_name = "rdma-agent"
}

locals {
  rdma_values = jsondecode(data.tencentcloud_tke_addon_values.rdma_agent.values)
}

# Read the effective memory limit of the agent container (default of the chart is 500Mi).
output "rdma_agent_memory_limit" {
  value = rdma_values.agent.resources.limits.memory
}
```
