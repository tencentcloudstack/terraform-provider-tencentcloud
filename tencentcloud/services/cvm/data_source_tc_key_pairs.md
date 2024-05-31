Use this data source to query key pairs.

Example Usage

Query key pairs by key ID

```hcl
data "tencentcloud_key_pairs" "key_id" {
  key_id = "skey-ie97i3ml"
}
```

Query key pairs by key name
```hcl
data "tencentcloud_key_pairs" "key_name" {
  key_name = "^test$"
}
```