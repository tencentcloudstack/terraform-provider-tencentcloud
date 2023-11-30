Use this data source to query detailed information of rum sign

Example Usage

```hcl
data "tencentcloud_rum_sign" "sign" {
  timeout   = 1800
  file_type = 1
}
```