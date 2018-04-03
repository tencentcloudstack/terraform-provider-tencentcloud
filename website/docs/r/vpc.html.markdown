---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc"
sidebar_current: "docs-tencentcloud-resource-vpc-x"
description: |-
  Provides an VPC resource.
---

# tencentcloud_vpc

Provides an VPC resource.

## Example Usage

Basic usage:

```hcl
resource "tencentcloud_vpc" "main" {
  name       = "my test vpc"
  cidr_block = "10.0.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for the VPC.
* `cidr_block` - (Required) The CIDR block for the VPC.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPC.
* `name` - The name for the VPC.
* `cidr_block` - The CIDR block of the VPC.
* `is_default` - Whether or not the default VPC.
* `is_multicast` - Whether or not the VPC has Multicast support.
