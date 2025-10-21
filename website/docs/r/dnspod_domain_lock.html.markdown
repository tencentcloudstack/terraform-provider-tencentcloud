---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_domain_lock"
sidebar_current: "docs-tencentcloud-resource-dnspod_domain_lock"
description: |-
  Provides a resource to create a dnspod domain_lock
---

# tencentcloud_dnspod_domain_lock

Provides a resource to create a dnspod domain_lock

## Example Usage

```hcl
resource "tencentcloud_dnspod_domain_lock" "domain_lock" {
  domain    = "dnspod.cn"
  lock_days = 30
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `lock_days` - (Required, Int, ForceNew) The number of max days to lock the domain+ Old packages: D_FREE 30 days, D_PLUS 90 days, D_EXTRA 30 days, D_EXPERT 60 days, D_ULTRA 365 days+ New packages: DP_FREE 365 days, DP_PLUS 365 days, DP_EXTRA 365 days, DP_EXPERT 365 days, DP_ULTRA 365 days.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `lock_code` - Domain unlock code, can be obtained through the ModifyDomainLock interface.


