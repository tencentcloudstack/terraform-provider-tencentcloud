---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule"
sidebar_current: "docs-tencentcloud-resource-waf_api_sec_sensitive_custom_api_exclude_rule"
description: |-
  Provides a resource to create a WAF api sec sensitive custom api exclude rule
---

# tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule

Provides a resource to create a WAF api sec sensitive custom api exclude rule

## Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule" "example" {
  domain     = "www.example.com"
  rule_name  = "tf-example"
  status     = 1
  match_type = "regex"
  content    = "/static"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `rule_name` - (Required, String, ForceNew) Rule name.
* `status` - (Required, Int) Rule switch, 0: off, 1: on.
* `content` - (Optional, String) Match content.
* `match_type` - (Optional, String) Match type, regex, prefix, suffix, contain match mode.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `update_time` - Update timestamp.


## Import

WAF api sec sensitive custom api exclude rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule.example www.example.com#tf-example
```

