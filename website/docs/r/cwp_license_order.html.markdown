---
subcategory: "Cwp"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cwp_license_order"
sidebar_current: "docs-tencentcloud-resource-cwp_license_order"
description: |-
  Provides a resource to create a CWP license order
---

# tencentcloud_cwp_license_order

Provides a resource to create a CWP license order

## Example Usage

```hcl
resource "tencentcloud_cwp_license_order" "example" {
  alias        = "tf_example"
  license_type = 0
  license_num  = 1
  region_id    = 1
  project_id   = 0
  tags = {
    createdBy = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `alias` - (Optional, String) Resource alias.
* `license_num` - (Optional, Int) Authorization quantity: the number of units that need to be purchased, The default is 1.
* `license_type` - (Optional, Int) Authorization type. 0: Pro Edition-pay-as-you-go; 1: Pro Edition-monthly subscription; 2 - Ultimate Edition-monthly subscriptionThe default is 0.
* `project_id` - (Optional, Int) Project ID. Default is 0.
* `region_id` - (Optional, Int) Region of purchase order. In this case, only 1 - Guangzhou and 9 - Singapore are supported. Guangzhou is recommended. Singapore region is reserved for allowlisted users. The default is 1.
* `tags` - (Optional, Map) Tags of the license order.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `license_id` - license id.
* `resource_id` - resource id.


## Import

CWP license order can be imported using the resourceId#regionId, e.g.

```
terraform import tencentcloud_cwp_license_order.example cwplic-130715d2#1
```

