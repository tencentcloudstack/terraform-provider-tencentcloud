---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_protection_mode"
sidebar_current: "docs-tencentcloud-resource-waf_protection_mode"
description: |-
  Provides a resource to create a waf protection_mode
---

# tencentcloud_waf_protection_mode

Provides a resource to create a waf protection_mode

## Example Usage

```hcl
resource "tencentcloud_waf_protection_mode" "example" {
  domain  = "keep.qcloudwaf.com"
  mode    = 10
  edition = "sparta-waf"
  type    = 0
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain.
* `mode` - (Required, Int) Protection status:10: Rule observation; AI off mode, 11: Rule observation; AI observation mode, 12: Rule observation; AI interception mode20: Rule interception; AI off mode, 21: Rule interception; AI observation mode, 22: Rule interception; AI interception mode.
* `edition` - (Optional, String) WAF edition. clb-waf means clb-waf, sparta-waf means saas-waf, default is sparta-waf.
* `type` - (Optional, Int) 0 is to modify the rule engine status, 1 is to modify the AI status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



