---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_cc"
sidebar_current: "docs-tencentcloud-resource-waf_cc"
description: |-
  Provides a resource to create a waf cc
---

# tencentcloud_waf_cc

Provides a resource to create a waf cc

## Example Usage

```hcl
resource "tencentcloud_waf_cc" "example" {
  domain      = "www.demo.com"
  name        = "terraform"
  status      = 1
  advance     = "0"
  limit       = "60"
  interval    = "60"
  url         = "/cc_demo"
  match_func  = 0
  action_type = "22"
  priority    = 50
  valid_time  = 600
  edition     = "sparta-waf"
  type        = 1
}
```

## Argument Reference

The following arguments are supported:

* `action_type` - (Required, String) Rule Action, 20 log, 21 captcha, 22 deny, 23 accurate deny.
* `advance` - (Required, String) Session match mode, 0 use session, 1 use ip.
* `domain` - (Required, String) Domain.
* `edition` - (Required, String) WAF edition.
* `interval` - (Required, String) Interval.
* `limit` - (Required, String) CC detection threshold.
* `match_func` - (Required, Int) Match method, 0 equal, 1 contains, 2 prefix.
* `name` - (Required, String) Rule Name.
* `priority` - (Required, Int) Rule Priority.
* `status` - (Required, Int) Rule Status, 0 rule close, 1 rule open.
* `url` - (Required, String) Check URL.
* `valid_time` - (Required, Int) Action ValidTime, minute unit.
* `event_id` - (Optional, String) Event ID.
* `session_applied` - (Optional, Set: [`Int`]) Advance mode use session id.
* `type` - (Optional, Int) Operate Type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


