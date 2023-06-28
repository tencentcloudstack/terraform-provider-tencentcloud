---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_proxy"
sidebar_current: "docs-tencentcloud-resource-cynosdb_proxy"
description: |-
  Provides a resource to create a cynosdb proxy
---

# tencentcloud_cynosdb_proxy

Provides a resource to create a cynosdb proxy

## Example Usage

```hcl
resource "tencentcloud_cynosdb_proxy" "proxy" {
  cluster_id               = "cynosdbmysql-bws8h88b"
  cpu                      = 2
  mem                      = 4000
  unique_vpc_id            = "vpc-k1t8ickr"
  unique_subnet_id         = "subnet-jdi5xn22"
  connection_pool_type     = "SessionConnectionPool"
  open_connection_pool     = "yes"
  connection_pool_time_out = 30
  security_group_ids       = ["sg-baxfiao5"]
  description              = "desc sample"
  proxy_zones {
    proxy_node_zone  = "ap-guangzhou-7"
    proxy_node_count = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `cpu` - (Required, Int) Number of CPU cores.
* `mem` - (Required, Int) Memory.
* `connection_pool_time_out` - (Optional, Int) Connection pool threshold: unit (second).
* `connection_pool_type` - (Optional, String) Connection pool type: SessionConnectionPool (session level Connection pool).
* `description` - (Optional, String) Description.
* `open_connection_pool` - (Optional, String) Whether to enable Connection pool, yes - enable, no - do not enable.
* `proxy_count` - (Optional, Int) Number of database proxy group nodes. If it is set at the same time as the `proxy_zones` field, the `proxy_zones` parameter shall prevail.
* `proxy_zones` - (Optional, List) Database node information.
* `security_group_ids` - (Optional, Set: [`String`]) Security Group ID Array.
* `unique_subnet_id` - (Optional, String) The private network subnet ID is consistent with the cluster subnet ID by default.
* `unique_vpc_id` - (Optional, String) Private network ID, which is consistent with the cluster private network ID by default.

The `proxy_zones` object supports the following:

* `proxy_node_count` - (Optional, Int) Number of proxy nodes.
* `proxy_node_zone` - (Optional, String) Proxy node availability zone.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `proxy_group_id` - Proxy Group Id.


