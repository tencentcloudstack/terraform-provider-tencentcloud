---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_instances_accept_attach"
sidebar_current: "docs-tencentcloud-resource-ccn_instances_accept_attach"
description: |-
  Provides a resource to create a vpc ccn_instances_accept_attach, you can use this resource to approve cross-region attachment.
---

# tencentcloud_ccn_instances_accept_attach

Provides a resource to create a vpc ccn_instances_accept_attach, you can use this resource to approve cross-region attachment.

## Example Usage

```hcl
resource "tencentcloud_ccn_instances_accept_attach" "ccn_instances_accept_attach" {
  ccn_id = "ccn-39lqkygf"
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
* `instances` - (Required, List, ForceNew) Accept List Of Attachment Instances.

The `instances` object supports the following:

* `instance_id` - (Required, String) Attachment Instance ID.
* `instance_region` - (Required, String) Instance Region.
* `description` - (Optional, String) Description.
* `instance_type` - (Optional, String) InstanceType: `VPC`, `DIRECTCONNECT`, `BMVPC`, `VPNGW`.
* `route_table_id` - (Optional, String) ID of the routing table associated with the instance. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



