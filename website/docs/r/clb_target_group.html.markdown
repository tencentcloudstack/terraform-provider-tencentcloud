---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_target_group"
sidebar_current: "docs-tencentcloud-resource-clb_target_group"
description: |-
  Provides a resource to create a CLB target group.
---

# tencentcloud_clb_target_group

Provides a resource to create a CLB target group.

## Example Usage

```hcl
resource "tencentcloud_clb_target_group" "test" {
  target_group_name = "test"
  port              = 33
}
```

## Argument Reference

The following arguments are supported:

* `port` - (Optional, Int) The default port of target group, add server after can use it.
* `target_group_instances` - (Optional, List, **Deprecated**) It has been deprecated from version 1.77.3. please use `tencentcloud_clb_target_group_instance_attachment` instead. The backend server of target group bind.
* `target_group_name` - (Optional, String) Target group name.
* `vpc_id` - (Optional, String, ForceNew) VPC ID, default is based on the network.

The `target_group_instances` object supports the following:

* `bind_ip` - (Required, String) The internal ip of target group instance.
* `port` - (Required, Int) The port of target group instance.
* `new_port` - (Optional, Int) The new port of target group instance.
* `weight` - (Optional, Int) The weight of target group instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CLB target group can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group.test lbtg-3k3io0i0
```

