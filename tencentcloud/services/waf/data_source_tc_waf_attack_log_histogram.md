Use this data source to query detailed information of waf attack_log_histogram

Example Usage

Obtain the specified domain name log information

```hcl
data "tencentcloud_waf_attack_log_histogram" "example" {
  domain       = "domain.com"
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-29 00:00:00"
  query_string = "method:GET"
}
```

Obtain all domain name log information

```hcl
data "tencentcloud_waf_attack_log_histogram" "example" {
  domain       = "all"
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-29 00:00:00"
  query_string = "method:GET"
}
```