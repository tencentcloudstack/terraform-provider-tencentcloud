---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_subnet"
sidebar_current: "docs-tencentcloud-resource-vpc-subnet"
description: |-
  Provides an Subnet resource.
---

# tencentcloud_subnet

Provides an VPC subnet resource.

## Example Usage

Basic usage:

```hcl
resource "tencentcloud_subnet" "main" {
  name              = "my test subnet"
  cidr_block        = "10.0.1.0/24"
  availability_zone = "ap-guangzhou-3"
  vpc_id            = "${tencent_vpc.main.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for the Subnet.
* `cidr_block` - (Required, Forces new resource) The CIDR block for the Subnet.
* `availability_zone`- (Required, Forces new resource) The AZ for the subnet.
* `vpc_id` - (Required, Forces new resource) The VPC ID.
* `route_table_id` - (Optional) The Route Table ID. Note that if this value is not set explicitly, the default route table ID will be used.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Subnet.
* `name` - The name for the Subnet.
* `cidr_block` - The CIDR block of the Subnet.
* `availability_zone`- The AZ for the subnet.
* `vpc_id` - The VPC ID.

