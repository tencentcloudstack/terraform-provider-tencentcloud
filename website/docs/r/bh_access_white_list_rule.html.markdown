---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_access_white_list_rule"
sidebar_current: "docs-tencentcloud-resource-bh_access_white_list_rule"
description: |-
  Provides a resource to create a BH access white list rule
---

# tencentcloud_bh_access_white_list_rule

Provides a resource to create a BH access white list rule

## Example Usage

```hcl
resource "tencentcloud_bh_access_white_list_rule" "example" {
  source = "1.1.1.1"
  remark = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `source` - (Required, String) IP address 10.10.10.1 or network segment 10.10.10.0/24, minimum length 4 bytes, maximum length 40 bytes.
* `remark` - (Optional, String) Remark information, minimum length 0 characters, maximum length 40 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - ID of the access white list rule.


## Import

BH access white list rule can be imported using the id, e.g.

```
terraform import tencentcloud_bh_access_white_list_rule.example 1235
```

