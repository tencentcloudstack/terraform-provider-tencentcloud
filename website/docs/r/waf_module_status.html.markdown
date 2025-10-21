---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_module_status"
sidebar_current: "docs-tencentcloud-resource-waf_module_status"
description: |-
  Provides a resource to create a waf module_status
---

# tencentcloud_waf_module_status

Provides a resource to create a waf module_status

## Example Usage

```hcl
resource "tencentcloud_waf_module_status" "example" {
  domain         = "demo.waf.com"
  web_security   = 1
  access_control = 0
  cc_protection  = 1
  api_protection = 1
  anti_tamper    = 1
  anti_leakage   = 0
}
```

## Argument Reference

The following arguments are supported:

* `access_control` - (Required, Int) ACL module status, 0:closed, 1:opened.
* `api_protection` - (Required, Int) API security module status, 0:closed, 1:opened.
* `cc_protection` - (Required, Int) CC module status, 0:closed, 1:opened.
* `domain` - (Required, String) Domain.
* `web_security` - (Required, Int) WEB security module status, 0:closed, 1:opened.
* `anti_leakage` - (Optional, Int) Anti leakage module status, 0:closed, 1:opened.
* `anti_tamper` - (Optional, Int) Anti tamper module status, 0:closed, 1:opened.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

waf module_status can be imported using the id, e.g.

```
terraform import tencentcloud_waf_module_status.example demo.waf.com
```

