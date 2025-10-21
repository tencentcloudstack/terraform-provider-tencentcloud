---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_attachment"
sidebar_current: "docs-tencentcloud-resource-ccn_attachment"
description: |-
  Provides a CCN attaching resource.
---

# tencentcloud_ccn_attachment

Provides a CCN attaching resource.

~> **NOTE:** The resource is no longer maintained, please use resource `tencentcloud_ccn_attachment_v2` instead.

## Example Usage

### Only Attachment instance

```hcl
variable "region" {
  default = "ap-guangzhou"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
}

variable "other_uin" {
  default = "100031344528"
}

variable "other_ccn" {
  default = "ccn-qhgojahx"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "172.16.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  is_multicast      = false
}

# create ccn
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  tags = {
    createBy = "terraform"
  }
}

# attachment instance
resource "tencentcloud_ccn_attachment" "attachment" {
  ccn_id          = tencentcloud_ccn.example.id
  instance_id     = tencentcloud_vpc.vpc.id
  instance_type   = "VPC"
  instance_region = var.region
}

# attachment other instance
resource "tencentcloud_ccn_attachment" "other_account" {
  ccn_id          = var.other_ccn
  instance_id     = tencentcloud_vpc.vpc.id
  instance_type   = "VPC"
  instance_region = var.region
  ccn_uin         = var.other_uin
}
```

### Attachment instance & route table

```hcl
variable "region" {
  default = "ap-guangzhou"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "172.16.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  is_multicast      = false
}

# create ccn
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  tags = {
    createBy = "terraform"
  }
}

# create ccn route table
resource "tencentcloud_ccn_route_table" "example" {
  ccn_id      = tencentcloud_ccn.example.id
  name        = "tf-example"
  description = "desc."
}

# attachment instance & route table
resource "tencentcloud_ccn_attachment" "attachment" {
  ccn_id          = tencentcloud_ccn.example.id
  instance_id     = tencentcloud_vpc.vpc.id
  instance_type   = "VPC"
  instance_region = var.region
  route_table_id  = tencentcloud_ccn_route_table.example.id
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String, ForceNew) ID of the CCN.
* `instance_id` - (Required, String, ForceNew) ID of instance is attached.
* `instance_region` - (Required, String, ForceNew) The region that the instance locates at.
* `instance_type` - (Required, String, ForceNew) Type of attached instance network, and available values include `VPC`, `DIRECTCONNECT`, `BMVPC` and `VPNGW`. Note: `VPNGW` type is only for whitelist customer now.
* `ccn_uin` - (Optional, String, ForceNew) Uin of the ccn attached. If not set, which means the uin of this account. This parameter is used with case when attaching ccn of other account to the instance of this account. For now only support instance type `VPC`.
* `description` - (Optional, String) Remark of attachment.
* `route_table_id` - (Optional, String, ForceNew) Ccn instance route table ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `attached_time` - Time of attaching.
* `cidr_block` - A network address block of the instance that is attached.
* `route_ids` - Route id list.
* `state` - States of instance is attached. Valid values: `PENDING`, `ACTIVE`, `EXPIRED`, `REJECTED`, `DELETED`, `FAILED`, `ATTACHING`, `DETACHING` and `DETACHFAILED`. `FAILED` means asynchronous forced disassociation after 2 hours. `DETACHFAILED` means asynchronous forced disassociation after 2 hours.


