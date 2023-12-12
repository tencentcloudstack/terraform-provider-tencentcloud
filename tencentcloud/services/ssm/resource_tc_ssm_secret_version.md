Provide a resource to create a SSM secret version.

-> **Note:** A maximum of 10 versions can be supported under one credential. Only new versions can be added to credentials in the enabled and disabled states.

Example Usage

Text type credential information plaintext

```hcl
resource "tencentcloud_ssm_secret" "example" {
  secret_name             = "tf-example"
  description             = "desc."
  recovery_window_in_days = 0
  is_enabled              = true

  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ssm_secret_version" "v1" {
  secret_name   = tencentcloud_ssm_secret.example.secret_name
  version_id    = "v1"
  secret_string = "this is secret string"
}
```

Binary credential information, encoded using base64

```hcl
resource "tencentcloud_ssm_secret_version" "v2" {
  secret_name   = tencentcloud_ssm_secret.example.secret_name
  version_id    = "v2"
  secret_binary = "MTIzMTIzMTIzMTIzMTIzQQ=="
}
```

Import

SSM secret version can be imported using the secretName#versionId, e.g.
```
$ terraform import tencentcloud_ssm_secret_version.v1 test#v1
```