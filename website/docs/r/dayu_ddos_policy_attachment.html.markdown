---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_ddos_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-dayu_ddos_policy_attachment"
description: |-
  Provides a resource to create a dayu DDoS policy attachment.
---

# tencentcloud_dayu_ddos_policy_attachment

Provides a resource to create a dayu DDoS policy attachment.

## Example Usage

```hcl
resource "tencentcloud_dayu_ddos_policy_attachment" "dayu_ddos_policy_attachment_basic" {
  resource_type = tencentcloud_dayu_ddos_policy.test_policy.resource_type
  resource_id   = "bgpip-00000294"
  policy_id     = tencentcloud_dayu_ddos_policy.test_policy.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew) Id of the policy.
* `resource_id` - (Required, ForceNew) Id of the attached resource.
* `resource_type` - (Required, ForceNew) Type of the resource that the DDoS policy works for, valid values are `bgpip`, `bgp`, `bgp-multip`, `net`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



