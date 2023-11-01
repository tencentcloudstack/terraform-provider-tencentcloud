---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_restart_nodes_operation"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_restart_nodes_operation"
description: |-
  Provides a resource to restart elasticsearch nodes
---

# tencentcloud_elasticsearch_restart_nodes_operation

Provides a resource to restart elasticsearch nodes

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_restart_nodes_operation" "restart_nodes_operation" {
  instance_id = "es-xxxxxx"
  node_names  = ["1648026612002990732"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `node_names` - (Required, Set: [`String`], ForceNew) List of node names.
* `force_restart` - (Optional, Bool, ForceNew) Whether to force a restart.
* `is_offline` - (Optional, Bool, ForceNew) Node status, used in blue-green mode; off-line node blue-green is risky.
* `restart_mode` - (Optional, String, ForceNew) Optional restart mode in-place,blue-green, which means restart and blue-green restart, respectively. The default is in-place.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



