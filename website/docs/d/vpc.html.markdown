---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc"
sidebar_current: "docs-tencentcloud-datasource-vpc-x"
description: |-
  Provides details about a specific VPC.
---

# tencentcloud_vpc

`tencentcloud_vpc` provides details about a specific VPC.

This resource can prove useful when a module accepts a vpc id as an input variable and needs to, for example, determine the CIDR block of that VPC.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_vpc_instances.

## Example Usage

The following example shows how one might accept a VPC id as a variable and use this data source to obtain the data necessary to create a subnet within it.

```hcl
variable "vpc_id" {}

data "tencentcloud_vpc" "selected" {
  id = "${var.vpc_id}"
}

resource "tencentcloud_subnet" "main" {
  name              = "my test subnet"
  cidr_block        = "${cidrsubnet(data.tencentcloud_vpc.selected.cidr_block, 4, 1)}"
  availability_zone = "eu-frankfurt-1"
  vpc_id            = "${data.tencentcloud_vpc.selected.id}"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The id of the specific VPC to retrieve.
* `name` - (Optional) VPC name. Fuzzy search is supported, as defined by [the underlying TencentCloud API](https://intl.cloud.tencent.com/document/product/215/1372).

## Attributes Reference

All of the argument attributes except `filter` blocks are also exported as result attributes. This data source will complete the data by populating any fields that are not included in the configuration with the data for the selected VPC.

The following attribute is additionally exported:

* `cidr_block` - The CIDR block of the VPC.
* `is_default` Whether or not the default VPC.
* `is_multicast` Whether or not the VPC has Multicast support.
