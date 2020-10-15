---
subcategory: "CLB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_targetgroup"
sidebar_current: "docs-tencentcloud-resource-clb_targetgroup"
description: |-
  Provides a resource to create a CLB target group.
---

# tencentcloud_clb_targetgroup

Provides a resource to create a CLB target group.

## Example Usage

```hcl
resource "tencentcloud_clb_targetgroup" "test" {
  target_group_name = "test"
  port              = 33
}
```

## Argument Reference

The following arguments are supported:

* `port` - (Optional) The default port for the target group.
* `target_group_name` - (Optional) Target group name.
* `vpc_id` - (Optional, ForceNew) VPC ID, default is based on the network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `target_group_id` - Target group ID.


## Import

CLB target group can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_targetgroup.test lbtg-3k3io0i0
```

