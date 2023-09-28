---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_custom_white_rule"
sidebar_current: "docs-tencentcloud-resource-waf_custom_white_rule"
description: |-
  Provides a resource to create a waf custom_white_rule
---

# tencentcloud_waf_custom_white_rule

Provides a resource to create a waf custom_white_rule

## Example Usage

```hcl
resource "tencentcloud_waf_custom_white_rule" "example" {
  name        = "tf-example"
  sort_id     = "30"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "1.1.1.1"
    arg          = ""
  }

  status = "1"
  domain = "test.com"
  bypass = "geoip,cc,owasp"
}
```

## Argument Reference

The following arguments are supported:

* `bypass` - (Required, String) Details of bypass.
* `domain` - (Required, String) Domain name that needs to add policy.
* `expire_time` - (Required, String) Expiration time, measured in seconds, such as 1677254399, which means the expiration time is 2023-02-24 23:59:59 0 means never expires.
* `name` - (Required, String) Rule Name.
* `sort_id` - (Required, String) Priority, value range 1-100, The smaller the number, the higher the execution priority of this rule.
* `strategies` - (Required, List) Strategies detail.
* `status` - (Optional, String) The status of the switch, 1 is on, 0 is off, default 1.

The `strategies` object supports the following:

* `arg` - (Required, String) Matching parameters.
* `compare_func` - (Required, String) Logical symbol.
* `content` - (Required, String) Matching Content.
* `field` - (Required, String) Matching Fields.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - rule ID.


## Import

waf custom_white_rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_custom_white_rule.example test.com#1100310837
```

