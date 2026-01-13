---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_package_instance"
sidebar_current: "docs-tencentcloud-resource-igtm_package_instance"
description: |-
  Provides a resource to create a IGTM package instance
---

# tencentcloud_igtm_package_instance

Provides a resource to create a IGTM package instance

## Example Usage

```hcl
resource "tencentcloud_igtm_package_instance" "example" {
  goods_type   = "STANDARD"
  auto_renew   = 1
  time_span    = 1
  auto_voucher = 1
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew` - (Required, Int) Auto renewal: 1 enable auto renewal; 2 disable auto renewal.
* `goods_type` - (Required, String) Package type: STANDARD for standard edition; ULTIMATE for flagship edition.
* `auto_voucher` - (Optional, Int) Whether to automatically select vouchers, 1 yes; 0 no, default is 0.
* `time_span` - (Optional, Int) Package duration in months, required for creation and renewal. Value range: 1~120.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `resource_id` - Resource ID.


## Import

IGTM package instance can be imported using the id, e.g.

```
terraform import tencentcloud_igtm_package_instance.example ins-wtqicjwzzze
```

