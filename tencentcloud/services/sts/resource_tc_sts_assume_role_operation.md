Provides a resource to perform an STS (Security Token Service) AssumeRole operation to obtain temporary security credentials.

Example Usage

Assume role with required parameters

```hcl
resource "tencentcloud_sts_assume_role_operation" "example" {
  role_arn          = "qcs::cam::uin/100000000001:roleName/testRoleName"
  role_session_name = "example-session"
}
```

Assume role with all parameters

```hcl
resource "tencentcloud_sts_assume_role_operation" "example" {
  role_arn          = "qcs::cam::uin/100000000001:roleName/testRoleName"
  role_session_name = "example-session"
  duration_seconds  = 7200
  policy            = jsonencode({
    version   = "2.0"
    statement = [
      {
        effect   = "allow"
        action   = ["cos:GetObject"]
        resource = ["*"]
      }
    ]
  })
  external_id     = "external-id-example"
  source_identity = "source-identity-example"
  serial_number   = "qcs::cam:uin/100000000001::mfa/softToken"
  token_code      = "123456"

  tags {
    key   = "env"
    value = "test"
  }

  tags {
    key   = "project"
    value = "demo"
  }
}
```
