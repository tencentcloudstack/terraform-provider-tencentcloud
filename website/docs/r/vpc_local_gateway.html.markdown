---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_local_gateway"
sidebar_current: "docs-tencentcloud-resource-vpc_local_gateway"
description: |-
  Provides a resource to create a vpc local_gateway
---

# tencentcloud_vpc_local_gateway

Provides a resource to create a vpc local_gateway

## Example Usage

```hcl
resource "tencentcloud_vpc_local_gateway" "local_gateway" {
  local_gateway_name = "local-gw-test"
  vpc_id             = "vpc-lh4nqig9"
  cdc_id             = "cluster-j9gyu1iy"
}
```

## Argument Reference

The following arguments are supported:

* `cdc_id` - (Required, String) CDC instance ID.
* `local_gateway_name` - (Required, String) Local gateway name.
* `vpc_id` - (Required, String) VPC instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc local_gateway can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_local_gateway.local_gateway local_gateway_id
```

