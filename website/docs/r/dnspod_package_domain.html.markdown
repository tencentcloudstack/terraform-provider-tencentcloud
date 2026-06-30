---
subcategory: "DNSPod"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_package_domain"
sidebar_current: "docs-tencentcloud-resource-dnspod_package_domain"
description: |-
  Provides a resource to manage DNSPod package domain binding
---

# tencentcloud_dnspod_package_domain

Provides a resource to manage DNSPod package domain binding

## Example Usage

```hcl
resource "tencentcloud_dnspod_package_domain" "package_domain" {
  resource_id = "res-xxxxx"
  domain_id   = 12345
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int) Domain ID to bind to the package.
* `resource_id` - (Required, String, ForceNew) Package resource ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `domain` - Domain.
* `downgrade` - Whether the package is downgraded.
* `grade_level` - Domain grade level.
* `grade_title` - Package grade title.
* `grade` - Package grade code.
* `is_grace_period` - Whether the package is in grace period.
* `remain_times` - Remaining domain bind/change times for the package.
* `status` - Package binding status.
* `vip_auto_renew` - VIP auto renew status. YES: enabled, NO: disabled, DEFAULT: default.
* `vip_end_at` - VIP end time.
* `vip_start_at` - VIP start time.


## Import

dnspod package_domain can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_package_domain.package_domain resource_id#domain_id
```

