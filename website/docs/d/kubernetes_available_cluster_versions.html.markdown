---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_available_cluster_versions"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_available_cluster_versions"
description: |-
  Use this data source to query detailed information of kubernetes available_cluster_versions
---

# tencentcloud_kubernetes_available_cluster_versions

Use this data source to query detailed information of kubernetes available_cluster_versions

## Example Usage

```hcl
data "tencentcloud_kubernetes_available_cluster_versions" "query_by_id" {
  cluster_id = "xxx"
}

output "versions_id" {
  description = "Query versions from id."
  value       = data.tencentcloud_kubernetes_available_cluster_versions.query_by_id.versions
}

data "tencentcloud_kubernetes_available_cluster_versions" "query_by_ids" {
  cluster_ids = ["xxx"]
}

output "versions_ids" {
  description = "Query versions from ids."
  value       = data.tencentcloud_kubernetes_available_cluster_versions.query_by_ids.clusters
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional, String) Cluster Id.
* `cluster_ids` - (Optional, Set: [`String`]) list of cluster IDs.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `clusters` - cluster information. Note: This field may return null, indicating that no valid value can be obtained.
  * `cluster_id` - Cluster ID.
  * `versions` - list of cluster major version numbers, for example 1.18.4.
* `versions` - Upgradable cluster version number. Note: This field may return null, indicating that no valid value can be obtained.


