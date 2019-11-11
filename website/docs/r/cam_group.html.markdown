---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_group"
sidebar_current: "docs-tencentcloud-resource-cam_group"
description: |-
  Provides a resource to create a CAM group.
---

# tencentcloud_cam_group

Provides a resource to create a CAM group.

## Example Usage

```hcl
resource "tencentcloud_cam_group" "foo" {
  name   = "cam-group-test"
  remark = "test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of CAM group.
* `remark` - (Optional) Description of the CAM group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Create time of the CAM group.


## Import

CAM group can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group.foo 90496
```

