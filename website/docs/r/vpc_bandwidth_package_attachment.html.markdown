---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_bandwidth_package_attachment"
sidebar_current: "docs-tencentcloud-resource-vpc_bandwidth_package_attachment"
description: |-
  Provides a resource to create a vpc bandwidth_package_attachment
---

# tencentcloud_vpc_bandwidth_package_attachment

Provides a resource to create a vpc bandwidth_package_attachment

## Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
}

resource "tencentcloud_vpc_bandwidth_package" "example" {
  network_type           = "BGP"
  charge_type            = "TOP5_POSTPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_clb_instance" "example" {
  network_type = "INTERNAL"
  clb_name     = "tf-example"
  project_id   = 0
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id

  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_vpc_bandwidth_package_attachment" "attachment" {
  resource_id          = tencentcloud_clb_instance.example.id
  bandwidth_package_id = tencentcloud_vpc_bandwidth_package.example.id
  network_type         = "BGP"
  resource_type        = "LoadBalance"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_package_id` - (Required, String, ForceNew) Bandwidth package unique ID, in the form of `bwp-xxxx`.
* `resource_id` - (Required, String, ForceNew) The unique ID of the resource, currently supports EIP resources and LB resources, such as `eip-xxxx`, `lb-xxxx`.
* `network_type` - (Optional, String, ForceNew) Bandwidth packet type, currently supports `BGP` type, indicating that the internal resource is BGP IP.
* `protocol` - (Optional, String, ForceNew) Bandwidth packet protocol type. Currently `ipv4` and `ipv6` protocol types are supported.
* `resource_type` - (Optional, String, ForceNew) Resource types, including `Address`, `LoadBalance`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



