---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_snapshot_policy_config"
sidebar_current: "docs-tencentcloud-resource-vpc_snapshot_policy_config"
description: |-
  Provides a resource to create a vpc snapshot_policy_config
---

# tencentcloud_vpc_snapshot_policy_config

Provides a resource to create a vpc snapshot_policy_config

## Example Usage

```hcl
resource "tencentcloud_vpc_snapshot_policy_config" "snapshot_policy_config" {
  snapshot_policy_id = "sspolicy-1t6cobbv"
  enable             = false
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Bool) If enable snapshot policy.
* `snapshot_policy_id` - (Required, String) Snapshot policy Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc snapshot_policy_config can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_snapshot_policy_config.snapshot_policy_config snapshot_policy_id
```

