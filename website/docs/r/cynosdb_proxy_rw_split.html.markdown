---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_proxy_rw_split"
sidebar_current: "docs-tencentcloud-resource-cynosdb_proxy_rw_split"
description: |-
  Provides a resource to create a cynosdb proxy_rw_split
---

# tencentcloud_cynosdb_proxy_rw_split

Provides a resource to create a cynosdb proxy_rw_split

## Example Usage

```hcl
resource "tencentcloud_cynosdb_proxy_rw_split" "proxy_rw_split" {
  cluster_id           = "cynosdbmysql-cgd2gpwr"
  proxy_group_id       = "cynosdbmysql-proxy-l6zf9t30"
  consistency_type     = "global"
  consistency_time_out = "30"
  weight_mode          = "system"
  instance_weights {
    instance_id = "cynosdbmysql-ins-9810be9i"
    weight      = 0
  }
  fail_over                = "yes"
  auto_add_ro              = "no"
  open_rw                  = "yes"
  rw_type                  = "READWRITE"
  trans_split              = false
  access_mode              = "balance"
  open_connection_pool     = "yes"
  connection_pool_type     = "SessionConnectionPool"
  connection_pool_time_out = 30
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `proxy_group_id` - (Required, String, ForceNew) Database Agent Group ID.
* `access_mode` - (Optional, String, ForceNew) Connection mode: nearby, balance.
* `auto_add_ro` - (Optional, String, ForceNew) Automatically add read-only instances, values: &amp;#39;yes&amp;#39;,&amp;#39; no &amp;#39;.
* `connection_pool_time_out` - (Optional, Int, ForceNew) Connection pool time.
* `connection_pool_type` - (Optional, String, ForceNew) Connection pool type: SessionConnectionPool.
* `consistency_time_out` - (Optional, String, ForceNew) Consistency timeout.
* `consistency_type` - (Optional, String, ForceNew) Consistency type; Eventual - Eventual consistency, session - session consistency, global - global consistency.
* `fail_over` - (Optional, String, ForceNew) Whether to enable failover. After the agent fails, the connection address will be routed to the main instance, with values of yes and no.
* `instance_weights` - (Optional, List, ForceNew) Instance read-only weight.
* `open_connection_pool` - (Optional, String, ForceNew) Open Connection pool: yes, no.
* `open_rw` - (Optional, String, ForceNew) Do you want to turn on read write separation.
* `rw_type` - (Optional, String, ForceNew) Read and write types: READWRITE, READONLY.
* `trans_split` - (Optional, Bool, ForceNew) Transaction splitting.
* `weight_mode` - (Optional, String, ForceNew) Reading and writing weight allocation mode; System automatic allocation: system, custom: custom.

The `instance_weights` object supports the following:

* `instance_id` - (Required, String) Instance Id.
* `weight` - (Required, Int) Instance Weight.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



