Provides a resource to create a CAM role.

Example Usage

Create normally

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

Or use service

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

Create with SAML provider

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

Import

CAM role can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_role.example 4611686018427733635
```