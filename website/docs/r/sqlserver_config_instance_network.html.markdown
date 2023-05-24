---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_instance_network"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_instance_network"
description: |-
  Provides a resource to create a sqlserver config_instance_network
---

# tencentcloud_sqlserver_config_instance_network

Provides a resource to create a sqlserver config_instance_network

## Example Usage

```hcl
resource "tencentcloud_sqlserver_config_instance_network" "config_instance_network" {
  instance_id   = "mssql-qelbzgwf"
  new_vpc_id    = "vpc-4owdpnwr"
  new_subnet_id = "sub-ahv6swf2"
  vip           = "172.16.16.48"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `new_subnet_id` - (Required, String) ID of the new subnet.
* `new_vpc_id` - (Required, String) ID of the new VPC.
* `vip` - (Optional, String) New VIP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver config_instance_network can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_instance_network.config_instance_network config_instance_network_id
```

