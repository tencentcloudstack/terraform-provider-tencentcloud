---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_js_injection_rule"
sidebar_current: "docs-tencentcloud-resource-teo_security_js_injection_rule"
description: |-
  Provides a resource to create a TEO security JavaScript injection rule.
---

# tencentcloud_teo_security_js_injection_rule

Provides a resource to create a TEO security JavaScript injection rule.

## Example Usage

```hcl
resource "tencentcloud_teo_security_js_injection_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"

  js_injection_rules {
    name       = "tf-example"
    priority   = 50
    condition  = "$${http.request.host} in ['www.demo.com']"
    inject_j_s = "inject-sdk-only"
  }
}
```

## Argument Reference

The following arguments are supported:

* `js_injection_rules` - (Required, List) JavaScript injection rule configuration. Only one rule is allowed per request.
* `zone_id` - (Required, String, ForceNew) Site ID.

The `js_injection_rules` object supports the following:

* `condition` - (Required, String) Match condition expression, e.g. `${http.request.host} in ['www.example.com'] and ${http.request.uri.path} in ['/path']`.
* `inject_j_s` - (Required, String) JavaScript injection option. Valid values: `no-injection` (do not inject JS); `inject-sdk-only` (inject SDK for all supported authentication methods, currently TC-RCE and TC-CAPTCHA).
* `name` - (Required, String) Rule name.
* `priority` - (Required, Int) Rule priority. Range: 0-100, smaller value means higher priority. Default: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TEO security JavaScript injection rule can be imported using the zoneId#ruleId, e.g.

```
terraform import tencentcloud_teo_security_js_injection_rule.example zone-3fkff38fyw8s#inject-0000040467
```

