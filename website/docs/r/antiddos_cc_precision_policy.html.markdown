---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_cc_precision_policy"
sidebar_current: "docs-tencentcloud-resource-antiddos_cc_precision_policy"
description: |-
  Provides a resource to create a antiddos cc_precision_policy
---

# tencentcloud_antiddos_cc_precision_policy

Provides a resource to create a antiddos cc_precision_policy

## Example Usage

```hcl
resource "tencentcloud_antiddos_cc_precision_policy" "cc_precision_policy" {
  instance_id   = "bgpip-0000078h"
  ip            = "212.64.62.191"
  protocol      = "http"
  domain        = "t.baidu.com"
  policy_action = "drop"
  policy_list {
    field_type     = "value"
    field_name     = "cgi"
    value          = "a.com"
    value_operator = "equal"
  }

  policy_list {
    field_type     = "value"
    field_name     = "ua"
    value          = "test"
    value_operator = "equal"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) domain.
* `instance_id` - (Required, String) Instance Id.
* `ip` - (Required, String) Ip value.
* `policy_action` - (Required, String) policy type, alg or drop.
* `policy_list` - (Required, List) policy list.
* `protocol` - (Required, String) protocol http or https.

The `policy_list` object supports the following:

* `field_name` - (Required, String) Configuration fields can take values of cgi, ua, cookie, referer, accept, srcip.
* `field_type` - (Required, String) field type.
* `value_operator` - (Required, String) Configuration item value comparison method, can take values of equal, not_ Equal, include.
* `value` - (Required, String) value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

antiddos cc_precision_policy can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_cc_precision_policy.cc_precision_policy ${instanceId}#${policyId}#${instanceIp}#${domain}#${protocol}
```

