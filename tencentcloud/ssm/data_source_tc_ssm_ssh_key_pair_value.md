Use this data source to query detailed information of ssm ssh_key_pair_value

~> **NOTE:** Must set at least one of `secret_name` or `ssh_key_id`.

Example Usage

```hcl
data "tencentcloud_ssm_ssh_key_pair_value" "example" {
  secret_name = "keep_terraform"
  ssh_key_id  = "skey-2ae2snwd"
}
```

Or

```hcl
data "tencentcloud_ssm_ssh_key_pair_value" "example" {
  secret_name = "keep_terraform"
}
```

Or

```hcl
data "tencentcloud_ssm_ssh_key_pair_value" "example" {
  ssh_key_id  = "skey-2ae2snwd"
}
```