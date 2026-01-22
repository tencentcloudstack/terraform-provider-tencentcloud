---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn"
sidebar_current: "docs-tencentcloud-resource-ccn"
description: |-
  Provides a resource to create a CCN instance.
---

# tencentcloud_ccn

Provides a resource to create a CCN instance.

~> **NOTE:** `route_overlap_flag` currently does not support setting to `false`.

## Example Usage

### Create a PREPAID CCN

```hcl
resource "tencentcloud_ccn" "example" {
  name                   = "tf-example"
  description            = "desc."
  qos                    = "AG"
  charge_type            = "PREPAID"
  bandwidth_limit_type   = "INTER_REGION_LIMIT"
  instance_metering_type = "BANDWIDTH"
  route_ecmp_flag        = true
  route_overlap_flag     = true

  tags = {
    createBy = "Terraform"
  }
}
```

### Create a POSTPAID regional export speed limit type CCN

```hcl
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "POSTPAID"
  bandwidth_limit_type = "OUTER_REGION_LIMIT"
  route_ecmp_flag      = false
  route_overlap_flag   = true
  tags = {
    createBy = "Terraform"
  }
}
```

### Create a POSTPAID inter-regional rate limit type CNN

```hcl
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "POSTPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the CCN to be queried, and maximum length does not exceed 60 bytes.
* `bandwidth_limit_type` - (Optional, String) The speed limit type. Valid values: `INTER_REGION_LIMIT`, `OUTER_REGION_LIMIT`. `OUTER_REGION_LIMIT` represents the regional export speed limit, `INTER_REGION_LIMIT` is the inter-regional speed limit. The default is `OUTER_REGION_LIMIT`.
* `charge_type` - (Optional, String, ForceNew) Billing mode. Valid values: `PREPAID`, `POSTPAID`. `PREPAID` means prepaid, which means annual and monthly subscription, `POSTPAID` means post-payment, which means billing by volume. The default is `POSTPAID`. The prepaid model only supports inter-regional speed limit, and the post-paid model supports inter-regional speed limit and regional export speed limit.
* `description` - (Optional, String) Description of CCN, and maximum length does not exceed 100 bytes.
* `instance_metering_type` - (Optional, String, ForceNew) Instance metering type. Valid values: `BANDWIDTH` (bandwidth billing), `TRAFFIC` (traffic billing). This parameter cannot be modified after creation.
* `qos` - (Optional, String, ForceNew) CCN service quality, 'PT': Platinum, 'AU': Gold, 'AG': Silver. The default is 'AU'.
* `route_ecmp_flag` - (Optional, Bool) Whether to enable the equivalent routing function. `true`: enabled, `false`: disabled. Default is false.
* `route_overlap_flag` - (Optional, Bool) Whether to enable the routing overlap function. `true`: enabled, `false`: disabled. Default is true, cannot set to false.
* `tags` - (Optional, Map) Instance tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of resource.
* `instance_count` - Number of attached instances.
* `state` - States of instance. Valid values: `ISOLATED`(arrears) and `AVAILABLE`.


## Import

Ccn instance can be imported, e.g.

```
$ terraform import tencentcloud_ccn.example ccn-al70jo89
```

