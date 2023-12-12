Use this data source to query detailed information of SSM secret version

Example Usage

```hcl
data "tencentcloud_ssm_secret_versions" "example" {
  secret_name = tencentcloud_ssm_secret_version.v1.secret_name
  version_id  = tencentcloud_ssm_secret_version.v1.version_id
}

resource "tencentcloud_ssm_secret" "example" {
  secret_name = "tf-example"
  description = "desc."

  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ssm_secret_version" "v1" {
  secret_name   = tencentcloud_ssm_secret.example.secret_name
  version_id    = "v1"
  secret_binary = "MTIzMTIzMTIzMTIzMTIzQQ=="
}
```