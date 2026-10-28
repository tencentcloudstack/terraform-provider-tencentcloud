---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_increase_plan_quota"
sidebar_current: "docs-tencentcloud-action-teo_increase_plan_quota"
description: |-
  Use this action to increase plan quota for TEO (EdgeOne) plans.
---

# tencentcloud_teo_increase_plan_quota

Use this action to increase plan quota for TEO (EdgeOne) plans.

This action calls the TEO `IncreasePlanQuota` API to purchase additional plan quotas (e.g., site count, precise access control rules, rate limiting rules). It is a one-time operation that does not persist any state.

## Example Usage

```hcl
action "tencentcloud_teo_increase_plan_quota" "example" {
  plan_id      = "edgeone-2unuvzjmmn2q"
  quota_type   = "site"
  quota_number = 10
}
```

## Argument Reference

The following arguments are supported:

* `plan_id` - (Required, String) Plan ID, in the format of edgeone-xxxxxxxx.
* `quota_number` - (Required, Int) Number of quotas to increase. Maximum 100 per request.
* `quota_type` - (Required, String) Quota type to increase. Valid values: `site`, `precise_access_control_rule`, `rate_limiting_rule`.
* `deal_name` - (Optional, String) Order number returned by the API.


