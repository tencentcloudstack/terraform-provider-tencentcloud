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

```hcl
resource "tencentcloud_vpc_resume_snapshot_instance" "resume_snapshot_instance" {
  snapshot_policy_id = "sspolicy-1t6cobbv"
  snapshot_file_id   = "ssfile-test"
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



