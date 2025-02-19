---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_custom_rule"
sidebar_current: "docs-tencentcloud-resource-waf_custom_rule"
description: |-
  Provides a resource to create a waf custom_rule
---

# tencentcloud_waf_custom_rule

Provides a resource to create a waf custom_rule

## Example Usage

```hcl
resource "tencentcloud_waf_custom_rule" "example" {
  name        = "tf-example"
  sort_id     = "50"
  redirect    = "/"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "2.2.2.2"
    arg          = ""
  }

  status      = "1"
  domain      = "test.com"
  action_type = "1"
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

waf custom_rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_custom_rule.example test.com#1100310609
```

