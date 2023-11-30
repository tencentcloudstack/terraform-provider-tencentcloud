Use this data source to query detailed information of waf attack_total_count

Example Usage

Obtain the specified domain name attack log

```hcl
data "tencentcloud_waf_attack_total_count" "example" {
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-07 00:00:00"
  domain       = "domain.com"
  query_string = "method:GET"
}
```

Obtain all domain name attack log

```hcl
data "tencentcloud_waf_attack_total_count" "example" {
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-07 00:00:00"
  domain       = "all"
  query_string = "method:GET"
}
```