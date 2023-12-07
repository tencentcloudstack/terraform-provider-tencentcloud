Use this data source to query key pairs.

Example Usage

```hcl
data "tencentcloud_key_pairs" "foo" {
  key_id = "skey-ie97i3ml"
}

data "tencentcloud_key_pairs" "name" {
  key_name = "^test$"
}
```