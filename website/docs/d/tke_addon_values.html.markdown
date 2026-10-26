---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tke_addon_values"
sidebar_current: "docs-tencentcloud-datasource-tke_addon_values"
description: |-
  Use this data source to query detailed information of TKE addon values.
---

# tencentcloud_tke_addon_values

Use this data source to query detailed information of TKE addon values.

-> **Note:** `values` is the rendering result of the parameters you configured merged with the
default parameters of the chart. It is **not** empty for an addon that is not installed in the
cluster, as long as the chart exists in the addon repository: in that case `values` is exactly
equal to `default_values` and contains no cluster information (for example `global.cluster.id`
is an empty string). To check whether an addon is actually installed, use the
`tencentcloud_kubernetes_addons` data source instead.

## Example Usage

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

### Example Usage of tuning an installed addon

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

## Argument Reference

The following arguments are supported:

* `addon_name` - (Required, String) Name of the addon, which can be obtained from the component management page in the cluster details page of the TKE console.
* `cluster_id` - (Required, String) ID of the TKE cluster, which can be obtained from the cluster list in the TKE console (https://console.cloud.tencent.com/tke2/cluster).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `default_values` - Parameter list supported by the addon in JSON format, whose value is the default value of the chart.
* `values` - Parameter list in JSON format. The value is the rendering result of the configured parameters merged with the default parameters of the chart. Note: it is NOT empty for an addon that is not installed in the cluster but whose chart exists in the addon repository - in that case the value is exactly the same as `default_values` and carries no cluster information (e.g. `global.cluster.id` is an empty string). Use the `tencentcloud_kubernetes_addons` data source to check whether the addon is installed.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeAddonValues](https://www.tencentcloud.com/document/api/457/56952)
