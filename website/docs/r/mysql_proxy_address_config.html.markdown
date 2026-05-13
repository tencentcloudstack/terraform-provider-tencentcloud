---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_proxy_address_config"
sidebar_current: "docs-tencentcloud-resource-mysql_proxy_address_config"
description: |-
  Provides a resource to manage MySQL database proxy address configuration.
---

# tencentcloud_mysql_proxy_address_config

Provides a resource to manage MySQL database proxy address configuration.

## Example Usage

```hcl
resource "tencentcloud_mysql_proxy_address_config" "example" {
  instance_id       = "cdb-o2t7gmjl"
  proxy_group_id    = "proxy-ov7dqp8n"
  proxy_address_id  = "proxyaddr-y8dnlfs0"
  weight_mode       = "system"
  is_kick_out       = true
  min_count         = 0
  max_delay         = 10
  fail_over         = true
  auto_add_ro       = true
  read_only         = false
  trans_split       = false
  connection_pool   = true
  auto_load_balance = true
  access_mode       = "nearby"
  proxy_allocation {
    region = "ap-guangzhou"
    zone   = "ap-guangzhou-6"

    proxy_instance {
      instance_id = "cdb-o2t7gmjl"
      weight      = 0
    }
  }

  proxy_allocation {
    region = "ap-guangzhou"
    zone   = "ap-guangzhou-7"

    proxy_instance {
      instance_id = "cdb-o2t7gmjl"
      weight      = 0
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `auto_add_ro` - (Required, Bool) Whether to automatically add read-only instances. Valid values: `true`, `false`.
* `fail_over` - (Required, Bool) Whether to enable failover. Valid values: `true`, `false`.
* `instance_id` - (Required, String, ForceNew) Instance ID, such as: cdb-xxxxxxxx.
* `is_kick_out` - (Required, Bool) Whether to enable delay elimination. Valid values: `true`, `false`.
* `max_delay` - (Required, Int) Delay elimination threshold in milliseconds. Value range: [1, 10000].
* `min_count` - (Required, Int) Minimum reserved quantity. Minimum value: 0. Note: only valid when IsKickOut is true.
* `proxy_address_id` - (Required, String, ForceNew) Proxy address ID, such as: proxyaddr-xxxxxxxx. Can be obtained through the DescribeCdbProxyInfo interface.
* `proxy_group_id` - (Required, String, ForceNew) Proxy group ID, such as: proxy-xxxxxxxx. Can be obtained through the DescribeCdbProxyInfo interface.
* `read_only` - (Required, Bool) Whether it is read-only. Valid values: `true`, `false`.
* `weight_mode` - (Required, String) Weight allocation mode. Valid values: `system` (system auto-allocation), `custom` (custom).
* `access_mode` - (Optional, String) Access mode. Valid values: `nearby` (nearby access), `balance` (load balancing). Default: `nearby`.
* `ap_node_as_ro_node` - (Optional, Bool) Whether to treat libra nodes as regular RO nodes.
* `ap_query_to_other_node` - (Optional, Bool) When libra node fails, whether to forward to other nodes.
* `auto_load_balance` - (Optional, Bool) Whether to enable adaptive load balancing. Default is disabled.
* `connection_pool` - (Optional, Bool) Whether to enable connection pool. Default is disabled. Note: for MySQL 8.0, the kernel minor version must be >= MySQL 8.0 20230630.
* `proxy_allocation` - (Optional, List) Read/write weight allocation. If WeightMode is `system`, the input weight does not take effect.
* `trans_split` - (Optional, Bool) Whether to enable transaction splitting. Default value: `false`.

The `proxy_allocation` object supports the following:

* `proxy_instance` - (Required, List) Proxy instance list.
* `region` - (Required, String) Region, such as: ap-guangzhou.
* `zone` - (Required, String) Availability zone, such as: ap-guangzhou-2.

The `proxy_instance` object of `proxy_allocation` supports the following:

* `instance_id` - (Required, String) Instance ID.
* `weight` - (Required, Int) Weight value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

MySQL proxy address config can be imported using the instanceId#proxyGroupId#proxyAddressId, e.g.

```
terraform import tencentcloud_mysql_proxy_address_config.example cdb-o2t7gmjl#proxy-ov7dqp8n#proxyaddr-y8dnlfs0
```

