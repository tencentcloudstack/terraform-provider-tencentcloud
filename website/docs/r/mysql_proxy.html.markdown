---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_proxy"
sidebar_current: "docs-tencentcloud-resource-mysql_proxy"
description: |-
  Provides a resource to create a mysql proxy
---

# tencentcloud_mysql_proxy

Provides a resource to create a mysql proxy

## Example Usage

```hcl
resource "tencentcloud_mysql_proxy" "proxy" {
  instance_id    = "cdb-fitq5t9h"
  uniq_vpc_id    = "vpc-4owdpnwr"
  uniq_subnet_id = "subnet-ahv6swf2"
  proxy_node_custom {
    node_count = 1
    cpu        = 2
    mem        = 4000
    region     = "ap-guangzhou"
    zone       = "ap-guangzhou-3"
  }
  security_group        = ["sg-edmur627"]
  desc                  = "desc1"
  connection_pool_limit = 2
  vip                   = "172.16.17.101"
  vport                 = 3306
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `proxy_node_custom` - (Required, List) Node specification configuration.
* `uniq_subnet_id` - (Required, String) Subnet id.
* `uniq_vpc_id` - (Required, String) Vpc id.
* `connection_pool_limit` - (Optional, Int) Connection Pool Threshold.
* `desc` - (Optional, String) Describe.
* `proxy_version` - (Optional, String) The current version of the database agent. No need to fill in when creating.
* `security_group` - (Optional, Set: [`String`]) Security group.
* `upgrade_time` - (Optional, String) Upgrade time: nowTime (upgrade completed) timeWindow (instance maintenance time), Required when modifying the agent version, No need to fill in when creating.
* `vip` - (Optional, String) IP address.
* `vport` - (Optional, Int) Port.

The `proxy_node_custom` object supports the following:

* `cpu` - (Required, Int) Number of CPU cores.
* `mem` - (Required, Int) Memory size.
* `node_count` - (Required, Int) Number of nodes.
* `region` - (Required, String) Region.
* `zone` - (Required, String) Zone.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `proxy_address_id` - Proxy address id.
* `proxy_group_id` - Proxy group id.


