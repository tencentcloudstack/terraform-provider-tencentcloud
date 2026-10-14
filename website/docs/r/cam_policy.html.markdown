---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_policy"
sidebar_current: "docs-tencentcloud-resource-cam_policy"
description: |-
  Provides a resource to create a CAM policy.
---

# tencentcloud_cam_policy

Provides a resource to create a CAM policy.

## Example Usage

```hcl
resource "tencentcloud_cam_policy" "foo" {
  name        = "tf_cam_policy"
  document    = <<EOF
{
  "version": "2.0",
  "statement": [
    {
      "action": [
        "name/sts:AssumeRole"
      ],
      "effect": "allow",
      "resource": [
        "*"
      ]
    }
  ]
}
EOF
  description = "tf_test"

  tags = {
    createBy = "Terraform"
    env      = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `document` - (Required, String) Document of the CAM policy. The syntax refers to [CAM POLICY](https://intl.cloud.tencent.com/document/product/598/10604). There are some notes when using this para in terraform: 1. The elements in JSON claimed supporting two types as `string` and `array` only support type `array`; 2. Terraform does not support the `root` syntax, when it appears, it must be replaced with the uin it stands for.
* `name` - (Required, String, ForceNew) Name of CAM policy.
* `description` - (Optional, String) Description of the CAM policy.
* `tags` - (Optional, Map) Instance tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the CAM policy.
* `type` - Type of the policy strategy. Valid values: `1`, `2`.  `1` means customer strategy and `2` means preset strategy.
* `update_time` - The last update time of the CAM policy.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreatePolicy](https://www.tencentcloud.com/document/product/598/32248) - This API (CreatePolicy) is used to create a policy.
* Create: [GetPolicy](https://www.tencentcloud.com/document/product/598/32244) - This API (GetPolicy) is used to query and view policy details.
* Read: [GetPolicy](https://www.tencentcloud.com/document/product/598/32244) - This API (GetPolicy) is used to query and view policy details.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Update: [UpdatePolicy](https://www.tencentcloud.com/document/product/598/45874) - This API is used to update a policy.
* Delete: [DeletePolicy](https://www.tencentcloud.com/document/product/598/32247) - This API (DeletePolicy) is used to delete a policy.

## Import

CAM policy can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_policy.foo 26655801
```

