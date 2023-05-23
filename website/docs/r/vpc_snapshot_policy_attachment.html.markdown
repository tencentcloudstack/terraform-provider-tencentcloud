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
resource "tencentcloud_vpc_snapshot_policy_attachment" "snapshot_policy_attachment" {
  snapshot_policy_id = "sspolicy-1t6cobbv"

  instances {
    instance_id     = "sg-r8ibzbd9"
    instance_name   = "cm-eks-cls-eizsc1iw-security-group"
    instance_region = "ap-guangzhou"
    instance_type   = "securitygroup"
  }
  instances {
    instance_id     = "sg-k3tn70lh"
    instance_name   = "keep-ci-temp-test-sg"
    instance_region = "ap-guangzhou"
    instance_type   = "securitygroup"
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

