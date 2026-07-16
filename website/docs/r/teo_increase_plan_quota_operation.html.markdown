---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_increase_plan_quota_operation"
sidebar_current: "docs-tencentcloud-resource-teo_increase_plan_quota_operation"
description: |-
  Provides a resource to increase plan quota for TEO (EdgeOne) plans.
---

# tencentcloud_teo_increase_plan_quota_operation

Provides a resource to increase plan quota for TEO (EdgeOne) plans.

## Example Usage

```hcl
resource "tencentcloud_teo_increase_plan_quota_operation" "example" {
  plan_id      = "edgeone-2unuvzjmmn2q"
  quota_type   = "site"
  quota_number = 10
}
```

## Argument Reference

The following arguments are supported:

* `plan_id` - (Required, String, ForceNew) Plan ID, in the format of edgeone-xxxxxxxx.
* `quota_number` - (Required, Int, ForceNew) Number of quotas to increase. Maximum 100 per request.
* `quota_type` - (Required, String, ForceNew) Quota type to increase. Valid values: site, precise_access_control_rule, rate_limiting_rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `deal_name` - Order number returned by the API.


