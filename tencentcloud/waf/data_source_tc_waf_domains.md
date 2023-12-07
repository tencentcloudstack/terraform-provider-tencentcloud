Use this data source to query detailed information of waf domains

Example Usage

Find all domains

```hcl
data "tencentcloud_waf_domains" "example" {}
```

Find domains by filter

```hcl
data "tencentcloud_waf_domains" "example" {
  instance_id = "waf_2kxtlbky01b3wceb"
  domain      = "tf.example.com"
}
```