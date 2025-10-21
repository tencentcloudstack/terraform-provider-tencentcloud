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
resource "tencentcloud_cos_bucket" "example" {
  bucket = "tf-example-1308919341"
  acl    = "private"
}

resource "tencentcloud_vpc_snapshot_policy" "example" {
  snapshot_policy_name = "tf-example"
  backup_type          = "time"
  cos_bucket           = tencentcloud_cos_bucket.example.bucket
  cos_region           = "ap-guangzhou"
  create_new_cos       = false
  keep_time            = 2

  backup_policies {
    backup_day  = "monday"
    backup_time = "00:00:00"
  }
  backup_policies {
    backup_day  = "tuesday"
    backup_time = "01:00:00"
  }
  backup_policies {
    backup_day  = "wednesday"
    backup_time = "02:00:00"
  }
}

resource "tencentcloud_vpc_snapshot_policy_config" "config" {
  snapshot_policy_id = tencentcloud_vpc_snapshot_policy.example.id
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

