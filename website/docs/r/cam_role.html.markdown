---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_role"
sidebar_current: "docs-tencentcloud-resource-cam_role"
description: |-
  Provides a resource to create a CAM role.
---

# tencentcloud_cam_role

Provides a resource to create a CAM role.

## Example Usage

Create normally

```hcl
resource "tencentcloud_cam_role" "foo" {
  name          = "cam-role-test"
  document      = <<EOF
{
  "version": "2.0",
  "statement": [
    {
      "action": ["name/sts:AssumeRole"],
      "effect": "allow",
      "principal": {
        "qcs": ["qcs::cam::uin/3374997817:uin/3374997817"]
      }
    }
  ]
}
EOF
  description   = "test"
  console_login = true
}
```

Create with SAML provider

```hcl
resource "tencentcloud_cam_role" "boo" {
  name          = "cam-role-test"
  document      = <<EOF
{
  "version": "2.0",
  "statement": [
    {
      "action": ["name/sts:AssumeRole", "name/sts:AssumeRoleWithWebIdentity"],
      "effect": "allow",
      "principal": {
        "federated": ["qcs::cam::uin/3374997817:saml-provider/XXXX-oooo"]
      }
    }
  ]
}
EOF
  description   = "test"
  console_login = true
}
```

## Argument Reference

The following arguments are supported:

* `document` - (Required) Document of the CAM role. The syntax refers to https://intl.cloud.tencent.com/document/product/598/10604. There are some notes when using this para in terraform: 1. The elements in json claimed supporting two types as `string` and `array` only support type `array`; 2. Terraform does not support the `root` syntax, when appears, it must be replaced with the uin it stands for.
* `name` - (Required, ForceNew) Name of CAM role.
* `console_login` - (Optional, ForceNew) Indicade whether the CAM role can login or not.
* `description` - (Optional) Description of the CAM role.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the CAM role.
* `update_time` - The last update time of the CAM role.


## Import

CAM role can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_role.foo 4611686018427733635
```

