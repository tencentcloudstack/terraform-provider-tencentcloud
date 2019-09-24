---
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
  ccn_id          = "${tencentcloud_ccn.main.id}"
  instance_type   = "VPC"
  instance_id     = "${tencentcloud_vpc.vpc.id}"
  instance_region = "${var.region}"
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, ForceNew) ID of the CCN.
* `instance_id` - (Required, ForceNew) ID of instance is attached.
* `instance_region` - (Required, ForceNew) The region that the instance locates at.
* `instance_type` - (Required, ForceNew) Type of attached instance network, and available values include VPC, DIRECTCONNECT and BMVPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `attached_time` - Time of attaching.
* `cidr_block` - A network address block of the instance that is attached.
* `state` - States of instance is attached, and available values include PENDING, ACTIVE, EXPIRED, REJECTED, DELETED, FAILED(asynchronous forced disassociation after 2 hours), ATTACHING, DETACHING and DETACHFAILED(asynchronous forced disassociation after 2 hours).


