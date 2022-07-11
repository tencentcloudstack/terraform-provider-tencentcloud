---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_levels"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_cluster_levels"
description: |-
  Provide a datasource to query TKE cluster levels.
---

# tencentcloud_kubernetes_cluster_levels

Provide a datasource to query TKE cluster levels.

## Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_levels" "foo" {}

output "level5" {
  value = data.tencentcloud_kubernetes_cluster_levels.foo.list.0.alias
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional, String) Specify cluster Id, if set will only query current cluster's available levels.
* `result_output_file` - (Optional, String) Used for save result.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of level information.
  * `alias` - Alias used for pass to cluster level arguments.
  * `config_map_count` - Number of ConfigMaps.
  * `crd_count` - Number of CRDs.
  * `enable` - Indicates whether the current level enabled.
  * `name` - Level name.
  * `node_count` - Number of nodes.
  * `other_count` - Number of others.
  * `pod_count` - Number of pods.


