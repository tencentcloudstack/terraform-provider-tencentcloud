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

### If type is v1

```hcl
resource "tencentcloud_clb_target_group" "example" {
  target_group_name = "tf-example"
  vpc_id            = "vpc-jy6pwoy2"
  port              = 8090
  type              = "v1"

  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

### If type is v2

```hcl
resource "tencentcloud_clb_target_group" "example" {
  target_group_name = "tf-example"
  vpc_id            = "vpc-jy6pwoy2"
  port              = 8090
  type              = "v2"
  protocol          = "TCP"
  weight            = 60

  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

### Or full_listen_switch is true

```hcl
resource "tencentcloud_clb_target_group" "example" {
  target_group_name  = "tf-example"
  vpc_id             = "vpc-jy6pwoy2"
  type               = "v2"
  protocol           = "TCP"
  weight             = 60
  full_listen_switch = true

  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `full_listen_switch` - (Optional, Bool) Full listening target group identifier, true indicates full listening target group, false indicates not full listening target group.
* `port` - (Optional, Int) The default port of target group, add server after can use it. If `full_listen_switch` is true, setting this parameter is not supported.
* `protocol` - (Optional, String) Target group backend forwarding protocol. This item is required for the v2 new version target group. Currently supports `TCP`, `UDP`.
* `tags` - (Optional, List) Label.
* `target_group_instances` - (Optional, List, **Deprecated**) It has been deprecated from version 1.77.3. please use `tencentcloud_clb_target_group_instance_attachment` instead. The backend server of target group bind.
* `target_group_name` - (Optional, String) Target group name.
* `type` - (Optional, String) Target group type, currently supports v1 (old version target group), v2 (new version target group), defaults to v1 (old version target group).
* `vpc_id` - (Optional, String, ForceNew) VPC ID, default is based on the network.
* `weight` - (Optional, Int) Default weights for backend services. Value range [0, 100]. After setting this value, when adding backend services to the target group, if the backend services do not have separate weights set, the default weights here will be used.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

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
$ terraform import tencentcloud_clb_target_group.example lbtg-3k3io0i0
```

