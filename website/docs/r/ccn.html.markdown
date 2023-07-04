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

## Example Usage

### Create a prepaid CCN

```hcl
resource "tencentcloud_ccn" "main" {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
}
```

### Create a post-paid regional export speed limit type CCN

```hcl
resource "tencentcloud_ccn" "main" {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
  qos                  = "AG"
  charge_type          = "POSTPAID"
  bandwidth_limit_type = "OUTER_REGION_LIMIT"
}
```

### Create a post-paid inter-regional rate limit type CNN

```hcl
resource "tencentcloud_ccn" "main" {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
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
* `qos` - (Optional, String, ForceNew) Service quality of CCN. Valid values: `PT`, `AU`, `AG`. The default is `AU`.
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
$ terraform import tencentcloud_ccn.test ccn-id
```

