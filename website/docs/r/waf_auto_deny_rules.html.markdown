---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_auto_deny_rules"
sidebar_current: "docs-tencentcloud-resource-waf_auto_deny_rules"
description: |-
  Provides a resource to create a waf auto_deny_rules
---

# tencentcloud_waf_auto_deny_rules

Provides a resource to create a waf auto_deny_rules

## Example Usage

```hcl
resource "tencentcloud_waf_auto_deny_rules" "example" {
  domain              = "demo.waf.com"
  attack_threshold    = 20
  time_threshold      = 12
  deny_time_threshold = 5
}
```

## Argument Reference

The following arguments are supported:

* `attack_threshold` - (Required, Int, ForceNew) The threshold number of attacks that triggers IP autodeny, ranging from 2 to 100 times.
* `deny_time_threshold` - (Required, Int, ForceNew) The IP autodeny time after triggering the IP autodeny, ranging from 5 to 360 minutes.
* `domain` - (Required, String, ForceNew) Domain.
* `time_threshold` - (Required, Int, ForceNew) IP autodeny statistical time, ranging from 1-60 minutes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

waf auto_deny_rules can be imported using the id, e.g.

```
terraform import tencentcloud_waf_auto_deny_rules.example demo.waf.com
```

