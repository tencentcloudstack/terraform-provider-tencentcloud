---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule"
sidebar_current: "docs-tencentcloud-resource-waf_api_sec_sensitive_custom_api_extract_rule"
description: |-
  Provides a resource to create a WAF api sec sensitive custom api extract rule
---

# tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule

Provides a resource to create a WAF api sec sensitive custom api extract rule

## Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule" "example" {
  domain    = "www.example.com"
  rule_name = "tf-example"
  status    = 1
  api_name  = "/api/login"
  methods   = ["GET", "POST"]
  regex     = "/api/.*"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `rule_name` - (Required, String, ForceNew) Rule name.
* `status` - (Required, Int) Rule switch, 0: off, 1: on.
* `api_name` - (Optional, String) API name.
* `methods` - (Optional, Set: [`String`]) Request method list.
* `regex` - (Optional, String) Regex match content.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `update_time` - Update timestamp.


## Import

WAF api sec sensitive custom api extract rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule.example www.example.com#tf-example
```

