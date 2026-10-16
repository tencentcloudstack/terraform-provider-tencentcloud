---
subcategory: "Cloud Access Management(CAM)"
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
  name   = "tf_cam_group"
  remark = "tf_group_remark"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of CAM group.
* `remark` - (Optional, String) Description of the CAM group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the CAM group.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateGroup](https://cloud.tencent.com/document/api/598/34582)
* Create: [GetGroup](https://cloud.tencent.com/document/api/598/34591)
* Read: [GetGroup](https://cloud.tencent.com/document/api/598/34591)
* Update: [UpdateGroup](https://cloud.tencent.com/document/api/598/34584)
* Delete: [DeleteGroup](https://cloud.tencent.com/document/api/598/34593)

## Import

CAM group can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group.foo 90496
```

