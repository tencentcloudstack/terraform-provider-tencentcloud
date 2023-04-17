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

## Example Usage

```hcl
variable "region" {
  default = "ap-guangzhou"
}

variable "otheruin" {
  default = "123353"
}

variable "otherccn" {
  default = "ccn-151ssaga"
}

resource "tencentcloud_vpc" "vpc" {
  name         = "ci-temp-test-vpc"
  cidr_block   = "10.0.0.0/16"
  dns_servers  = ["119.29.29.29", "8.8.8.8"]
  is_multicast = false
}

resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_ccn_attachment" "attachment" {
  ccn_id          = tencentcloud_ccn.main.id
  instance_type   = "VPC"
  instance_id     = tencentcloud_vpc.vpc.id
  instance_region = var.region
}

resource "tencentcloud_ccn_attachment" "other_account" {
  ccn_id          = var.otherccn
  instance_type   = "VPC"
  instance_id     = tencentcloud_vpc.vpc.id
  instance_region = var.region
  ccn_uin         = var.otheruin
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String, ForceNew) ID of the CCN.
* `instance_id` - (Required, String, ForceNew) ID of instance is attached.
* `instance_region` - (Required, String, ForceNew) The region that the instance locates at.
* `instance_type` - (Required, String, ForceNew) Type of attached instance network, and available values include `VPC`, `DIRECTCONNECT`, `BMVPC` and `VPNGW`. Note: `VPNGW` type is only for whitelist customer now.
* `ccn_uin` - (Optional, String, ForceNew) Uin of the ccn attached. Default is ``, which means the uin of this account. This parameter is used with case when attaching ccn of other account to the instance of this account. For now only support instance type `VPC`.
* `description` - (Optional, String) Remark of attachment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `attached_time` - Time of attaching.
* `cidr_block` - A network address block of the instance that is attached.
* `state` - States of instance is attached. Valid values: `PENDING`, `ACTIVE`, `EXPIRED`, `REJECTED`, `DELETED`, `FAILED`, `ATTACHING`, `DETACHING` and `DETACHFAILED`. `FAILED` means asynchronous forced disassociation after 2 hours. `DETACHFAILED` means asynchronous forced disassociation after 2 hours.


