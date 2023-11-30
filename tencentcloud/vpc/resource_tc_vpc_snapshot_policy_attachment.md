Provides a resource to create a vpc snapshot_policy_attachment

Example Usage

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
    instance_type      = "securitygroup"
    instance_id        = tencentcloud_security_group.example.id
    instance_name      = "tf-example"
    instance_region    = "ap-guangzhou"
  }
}
```

Import

vpc snapshot_policy_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_snapshot_policy_attachment.snapshot_policy_attachment snapshot_policy_attachment_id
```