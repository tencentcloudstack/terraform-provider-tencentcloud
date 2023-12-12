Use this data source to query detailed information of ssm rotation_history

Example Usage

```hcl
data "tencentcloud_ssm_rotation_history" "example" {
  secret_name = "keep_terraform"
}
```