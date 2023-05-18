---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_traffic_package"
sidebar_current: "docs-tencentcloud-resource-vpc_traffic_package"
description: |-
  Provides a resource to create a vpc traffic_package
---

# tencentcloud_vpc_traffic_package

Provides a resource to create a vpc traffic_package

## Example Usage

```hcl
resource "tencentcloud_vpc_traffic_package" "traffic_package" {
  traffic_amount = 10
}
```

## Argument Reference

The following arguments are supported:

* `traffic_amount` - (Required, Int, ForceNew) Traffic Package Amount, eg: 10,20,50,512,1024,5120,51200,60,300,600,3072,6144,30720,61440,307200.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Created time.
* `remaining_amount` - Remaining amount.
* `used_amount` - Used amount.


## Import

vpc traffic_package can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_traffic_package.traffic_package traffic_package_id
```

