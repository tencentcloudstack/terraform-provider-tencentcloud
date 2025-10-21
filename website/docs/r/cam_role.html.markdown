---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_role"
sidebar_current: "docs-tencentcloud-resource-cam_role"
description: |-
  Provides a resource to create a CAM role.
---

# tencentcloud_cam_role

Provides a resource to create a CAM role.

## Example Usage

### Create normally

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  uin = data.tencentcloud_user_info.info.owner_uin
}

resource "tencentcloud_cam_role" "example" {
  name = "tf-example"
  document = jsonencode(
    {
      statement = [
        {
          action = "name/sts:AssumeRole"
          effect = "allow"
          principal = {
            qcs = [
              "qcs::cam::uin/${local.uin}:root",
            ]
          }
        },
      ]
      version = "2.0"
    }
  )
  console_login    = true
  description      = "test"
  session_duration = 7200
  tags = {
    createBy = "Terraform"
  }
}

output "uin" {
  value = local.uin
}

output "arn" {
  value = tencentcloud_cam_role.example.role_arn
}
```

### Or use service

```hcl
resource "tencentcloud_cam_role" "example" {
  name = "tf-example"
  document = jsonencode(
    {
      statement = [
        {
          action = "name/sts:AssumeRole"
          effect = "allow"
          principal = {
            service = [
              "scf.qcloud.com",
            ]
          }
        },
      ]
      version = "2.0"
    }
  )
  console_login    = true
  description      = "test"
  session_duration = 7200
  tags = {
    createBy = "Terraform"
  }
}
```

### Create with SAML provider

```hcl
variable "saml-provider" {
  default = "example"
}

data "tencentcloud_user_info" "info" {}

locals {
  uin           = data.tencentcloud_user_info.info.uin
  saml_provider = var.saml-provider
}

resource "tencentcloud_cam_role" "example" {
  name          = "tf-example"
  document      = <<EOF
{
  "version": "2.0",
  "statement": [
    {
      "action": [
        "name/sts:AssumeRole"
      ],
      "effect": "allow",
      "principal": {
        "qcs": [
          "qcs::cam::uin/${local.uin}:saml-provider/${local.saml_provider}"
        ]
      }
    }
  ]
}
EOF
  description   = "terraform demo"
  console_login = true
}
```

## Argument Reference

The following arguments are supported:

* `document` - (Required, String) Document of the CAM role. The syntax refers to [CAM POLICY](https://intl.cloud.tencent.com/document/product/598/10604). There are some notes when using this para in terraform: 1. The elements in json claimed supporting two types as `string` and `array` only support type `array`; 2. Terraform does not support the `root` syntax, when appears, it must be replaced with the uin it stands for.
* `name` - (Required, String, ForceNew) Name of CAM role.
* `console_login` - (Optional, Bool) Indicates whether the CAM role can login or not.
* `description` - (Optional, String) Description of the CAM role.
* `session_duration` - (Optional, Int) The maximum validity period of the temporary key for creating a role.
* `tags` - (Optional, Map) A list of tags used to associate different resources.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the CAM role.
* `role_arn` - RoleArn Information for Roles.
* `update_time` - The last update time of the CAM role.


## Import

CAM role can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_role.example 4611686018427733635
```

