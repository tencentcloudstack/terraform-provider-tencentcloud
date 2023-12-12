Use this data source to query detailed information of SSM secret

Example Usage

```hcl
data "tencentcloud_ssm_secrets" "example" {
  secret_name = tencentcloud_ssm_secret.example.secret_name
  state       = 1
}

resource "tencentcloud_ssm_secret" "example" {
  secret_name = "tf_example"
  description = "desc."

  tags = {
    createdBy = "terraform"
  }
}
```

OR you can filter by tags

```hcl
data "tencentcloud_ssm_secrets" "example" {
  secret_name = tencentcloud_ssm_secret.example.secret_name
  state       = 1

  tags = {
    createdBy = "terraform"
  }
}
```