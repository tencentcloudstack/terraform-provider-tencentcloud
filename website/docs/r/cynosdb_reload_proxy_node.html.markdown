---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_reload_proxy_node"
sidebar_current: "docs-tencentcloud-resource-cynosdb_reload_proxy_node"
description: |-
  Provides a resource to create a cynosdb reload_proxy_node
---

# tencentcloud_cynosdb_reload_proxy_node

Provides a resource to create a cynosdb reload_proxy_node

## Example Usage

```hcl
resource "tencentcloud_cynosdb_reload_proxy_node" "reload_proxy_node" {
  cluster_id     = "cynosdbmysql-cgd2gpwr"
  proxy_group_id = "cynosdbmysql-proxy-8lqtl8pk"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) cluster id.
* `proxy_group_id` - (Required, String, ForceNew) proxy group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cynosdb reload_proxy_node can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_reload_proxy_node.reload_proxy_node reload_proxy_node_id
```

