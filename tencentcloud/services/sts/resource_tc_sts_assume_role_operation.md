Provides a resource to perform an STS (Security Token Service) AssumeRole operation to obtain temporary security credentials.

Example Usage

Assume role with required parameters

```hcl
resource "tencentcloud_sts_assume_role_operation" "example" {
  role_arn          = "qcs::cam::uin/100000000000:roleName/tf-example"
  role_session_name = "tf-example"
}
```
