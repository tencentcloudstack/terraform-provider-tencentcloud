---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_resume_snapshot_instance"
sidebar_current: "docs-tencentcloud-resource-vpc_resume_snapshot_instance"
description: |-
  Provides a resource to create a vpc resume_snapshot_instance
---

# tencentcloud_vpc_resume_snapshot_instance

Provides a resource to create a vpc resume_snapshot_instance

## Example Usage

### Basic example

```hcl
resource "tencentcloud_vpc_resume_snapshot_instance" "resume_snapshot_instance" {
  snapshot_policy_id = "sspolicy-1t6cobbv"
  snapshot_file_id   = "ssfile-emtabuwu2z"
  instance_id        = "ntrgm89v"
}
```

### Complete example

```hcl
data "tencentcloud_vpc_snapshot_files" "example" {
  business_type = "securitygroup"
  instance_id   = "sg-902tl7t7"
  start_date    = "2022-10-10 00:00:00"
  end_date      = "2023-10-30 00:00:00"
}

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

resource "tencentcloud_vpc_resume_snapshot_instance" "example" {
  snapshot_policy_id = tencentcloud_vpc_snapshot_policy.example.id
  snapshot_file_id   = data.tencentcloud_vpc_snapshot_files.example.snapshot_file_set.0.snapshot_file_id
  instance_id        = "policy-1t6cob"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) InstanceId.
* `snapshot_file_id` - (Required, String, ForceNew) Snapshot file Id.
* `snapshot_policy_id` - (Required, String, ForceNew) Snapshot policy Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



