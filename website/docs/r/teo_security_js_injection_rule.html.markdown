---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_js_injection_rule"
sidebar_current: "docs-tencentcloud-resource-teo_security_js_injection_rule"
description: |-
  Provides a resource to create a TEO security js injection rule
---

# tencentcloud_teo_security_js_injection_rule

Provides a resource to create a TEO security js injection rule

## Example Usage

```hcl
resource "tencentcloud_teo_security_js_injection_rule" "example" {
  zone_id = "zone-2qtuhspy7cr6"
  js_injection_rules {
    name      = "test-rule-1"
    condition = "${http.request.host} in ['example.com']"
    inject_js = "inject-sdk-only"
  }
  js_injection_rules {
    name      = "test-rule-2"
    priority  = 10
    condition = "${http.request.uri.path} in ['/api/*']"
    inject_js = "no-injection"
  }
}
```

## Argument Reference

The following arguments are supported:

* `js_injection_rules` - (Required, List) JavaScript injection rule list.
* `zone_id` - (Required, String, ForceNew) Site ID.

The `js_injection_rules` object supports the following:

* `condition` - (Required, String) Match condition content. Must conform to expression syntax.
* `name` - (Required, String) Rule name.
* `inject_js` - (Optional, String) JavaScript injection option. Valid values: `no-injection`, `inject-sdk-only`.
* `priority` - (Optional, Int) Rule priority. The smaller the value, the earlier the execution. Range 0-100, default 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `js_injection_rule_ids` - Rule ID list.


## Import

TEO security js injection rule can be imported using the zone_id, e.g.

```
terraform import tencentcloud_teo_security_js_injection_rule.example zone-2qtuhspy7cr6
```

