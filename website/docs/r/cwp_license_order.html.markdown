---
subcategory: "Cwp"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cwp_license_order"
sidebar_current: "docs-tencentcloud-resource-cwp_license_order"
description: |-
  Provides a resource to create a cwp license_order
---

# tencentcloud_cwp_license_order

Provides a resource to create a cwp license_order

## Example Usage

```hcl
resource "tencentcloud_cwp_license_order" "example" {
  alias        = "tf_example"
  license_type = 0
  license_num  = 1
  region_id    = 1
  project_id   = 0
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `alias` - (Optional, String) Resource alias.
* `license_num` - (Optional, Int) License quantity, Quantity to be purchased.Default is 1.
* `license_type` - (Optional, Int) LicenseType, 0 CWP Pro - Pay as you go, 1 CWP Pro - Monthly subscription, 2 CWP Ultimate - Monthly subscription. Default is 0.
* `project_id` - (Optional, Int) Project ID. Default is 0.
* `region_id` - (Optional, Int) Purchase order region, only 1 Guangzhou, 9 Singapore is supported here. Guangzhou is recommended. Singapore is whitelisted. Default is 1.
* `tags` - (Optional, Map) Tags of the license order.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `license_id` - license id.
* `resource_id` - resource id.


## Import

cwp license_order can be imported using the id, e.g.

```
terraform import tencentcloud_cwp_license_order.example cwplic-130715d2#1
```

