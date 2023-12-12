Use this data source to query detailed information of waf find_domains

Example Usage

Find all domains

```hcl
data "tencentcloud_waf_find_domains" "example" {}
```

Find domains by filter

```hcl
data "tencentcloud_waf_find_domains" "example" {
  key           = "keyWord"
  is_waf_domain = "1"
  by            = "FindTime"
  order         = "asc"
}
```