Provides a resource to create a CAM role.

Example Usage

Create normally

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  uin = data.tencentcloud_user_info.info.owner_uin
}

output "uin" {
  value = local.uin
}

resource "tencentcloud_cam_role" "foo" {
  name     = "cam-role-test"
  document = jsonencode(
    {
      statement = [
        {
          action    = "name/sts:AssumeRole"
          effect    = "allow"
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
  tags             = {
    test  = "tf-cam-role"
  }
}
```

Create with SAML provider

```hcl
variable "saml-provider" {
  default = "example"
}

locals {
  uin = data.tencentcloud_user_info.info.uin
  saml_provider = var.saml-provider
}

data "tencentcloud_user_info" "info" {}

resource "tencentcloud_cam_role" "boo" {
  name          = "tf_cam_role"
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
  description   = "tf_test"
  console_login = true
}
```

Import

CAM role can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_role.foo 4611686018427733635
```