---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_snapshot_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-vpc_snapshot_policy_attachment"
description: |-
  Provides a resource to create a vpc snapshot_policy_attachment
---

# tencentcloud_vpc_snapshot_policy_attachment

Provides a resource to create a vpc snapshot_policy_attachment

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

resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "desc."
}

resource "tencentcloud_vpc_snapshot_policy_attachment" "attachment" {
  snapshot_policy_id = tencentcloud_vpc_snapshot_policy.example.id

  instances {
    instance_type   = "securitygroup"
    instance_id     = tencentcloud_security_group.example.id
    instance_name   = "tf-example"
    instance_region = "ap-guangzhou"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instances` - (Required, Set, ForceNew) Associated instance information.
* `snapshot_policy_id` - (Required, String, ForceNew) Snapshot policy Id.

The `instances` object supports the following:

* `instance_id` - (Required, String, ForceNew) InstanceId.
* `instance_region` - (Required, String, ForceNew) The region where the instance is located.
* `instance_type` - (Required, String, ForceNew) Instance type, currently supports set: `securitygroup`.
* `instance_name` - (Optional, String, ForceNew) Instance name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc snapshot_policy_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_snapshot_policy_attachment.snapshot_policy_attachment snapshot_policy_attachment_id
```

