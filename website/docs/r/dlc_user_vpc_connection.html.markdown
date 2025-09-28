---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_user_vpc_connection"
sidebar_current: "docs-tencentcloud-resource-dlc_user_vpc_connection"
description: |-
  Provides a resource to create a DLC user vpc connection
---

# tencentcloud_dlc_user_vpc_connection

Provides a resource to create a DLC user vpc connection

## Example Usage

```hcl
resource "tencentcloud_dlc_user_vpc_connection" "example" {
  user_vpc_id            = "vpc-f7fa1fu5"
  user_subnet_id         = "subnet-ds2t3udw"
  user_vpc_endpoint_name = "tf-example"
  engine_network_id      = "DataEngine-Network-2mfg9icb"
}
```

### Or

```hcl
resource "tencentcloud_dlc_user_vpc_connection" "example" {
  user_vpc_id            = "vpc-f7fa1fu5"
  user_subnet_id         = "subnet-ds2t3udw"
  user_vpc_endpoint_name = "tf-example"
  engine_network_id      = "DataEngine-Network-2mfg9icb"
  user_vpc_endpoint_vip  = "10.0.1.10"
}
```

## Argument Reference

The following arguments are supported:

* `engine_network_id` - (Required, String, ForceNew) Engine network ID.
* `user_subnet_id` - (Required, String, ForceNew) User subnet ID.
* `user_vpc_endpoint_name` - (Required, String, ForceNew) User vpc endpoint name.
* `user_vpc_id` - (Required, String, ForceNew) User vpc ID.
* `user_vpc_endpoint_vip` - (Optional, String, ForceNew) Manually specify VIP, if not filled in, an IP address under the subnet will be automatically assigned.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `user_vpc_endpoint_id` - User endpoint ID.


