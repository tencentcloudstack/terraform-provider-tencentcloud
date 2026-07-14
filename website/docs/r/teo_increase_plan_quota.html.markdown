---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_increase_plan_quota"
sidebar_current: "docs-tencentcloud-resource-teo_increase_plan_quota"
description: |-
  Provides a resource to increase TEO plan quota. Use this resource to purchase additional quota for a TEO enterprise plan when the number of bound sites, Web Protection custom precise match policy rules, or rate limiting precise rate limiting module rules reaches the plan's quota limit.
---

# tencentcloud_teo_increase_plan_quota

Provides a resource to increase TEO plan quota. Use this resource to purchase additional quota for a TEO enterprise plan when the number of bound sites, Web Protection custom precise match policy rules, or rate limiting precise rate limiting module rules reaches the plan's quota limit.

## Example Usage

```hcl
resource "tencentcloud_teo_increase_plan_quota" "example" {
  plan_id      = "edgeone-2unuvzjmmn2q"
  quota_type   = "site"
  quota_number = 1
}
```

## Argument Reference

The following arguments are supported:

* `plan_id` - (Required, String, ForceNew) Plan ID, e.g., edgeone-2unuvzjmmn2q.
* `quota_number` - (Required, Int, ForceNew) Number of quotas to increase. Maximum is 100 per request.
* `quota_type` - (Required, String, ForceNew) Quota type. Valid values: `site` (site count), `precise_access_control_rule` (Web Protection - Custom Rules - Precise Match Policy rule quota), `rate_limiting_rule` (Web Protection - Rate Limiting - Precise Rate Limiting module rule quota).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `deal_name` - Order number returned after successful quota increase.


