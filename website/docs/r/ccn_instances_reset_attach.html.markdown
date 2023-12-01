---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_instances_reset_attach"
sidebar_current: "docs-tencentcloud-resource-ccn_instances_reset_attach"
description: |-
  Provides a resource to create a vpc ccn_instances_reset_attach, you can use this resource to reset cross-region attachment.
---

# tencentcloud_ccn_instances_reset_attach

Provides a resource to create a vpc ccn_instances_reset_attach, you can use this resource to reset cross-region attachment.

## Example Usage

```hcl
resource "tencentcloud_ccn_instances_reset_attach" "ccn_instances_reset_attach" {
  ccn_id  = "ccn-39lqkygf"
  ccn_uin = "100022975249"
  instances {
    instance_id     = "vpc-j9yhbzpn"
    instance_region = "ap-guangzhou"
    instance_type   = "VPC"
  }
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String, ForceNew) CCN Instance ID.
* `ccn_uin` - (Required, String, ForceNew) CCN Uin (root account).
* `instances` - (Required, List, ForceNew) List Of Attachment Instances.

The `instances` object supports the following:

* `instance_id` - (Required, String) Attachment Instance ID.
* `instance_region` - (Required, String) Instance Region.
* `description` - (Optional, String) Description.
* `instance_type` - (Optional, String) InstanceType: `VPC`, `DIRECTCONNECT`, `BMVPC`, `VPNGW`.
* `route_table_id` - (Optional, String) ID of the routing table associated with the instance. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



