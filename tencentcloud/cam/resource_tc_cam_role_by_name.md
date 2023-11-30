Provides a resource to create a CAM role.

Example Usage

Create normally

```hcl
resource "tencentcloud_cam_role_by_name" "foo" {
  name          = "tf_cam_role"
  document      = <<EOF
{
  "version": "2.0",
  "statement": [
    {
      "action": ["name/sts:AssumeRole"],
      "effect": "allow",
      "principal": {
        "qcs": ["qcs::cam::uin/<your-account-id>:uin/<your-account-id>"]
      }
    }
  ]
}
EOF
  description   = "test"
  console_login = true
  tags = {
    test  = "tf-cam-role",
  }
}
```

Create with SAML provider

```hcl
resource "tencentcloud_cam_role_by_name" "boo" {
  name          = "cam-role-test"
  document      = <<EOF
{
  "version": "2.0",
  "statement": [
    {
      "action": ["name/sts:AssumeRole", "name/sts:AssumeRoleWithWebIdentity"],
      "effect": "allow",
      "principal": {
        "federated": ["qcs::cam::uin/<your-account-id>:saml-provider/<your-name>"]
      }
    }
  ]
}
EOF
  description   = "test"
  console_login = true
}
```

Import

CAM role can be imported using the name, e.g.

```
$ terraform import tencentcloud_cam_role_by_name.foo cam-role-test
```