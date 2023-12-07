Use this data source to query detailed information of waf attack_log_list

Example Usage

Obtain the specified domain name attack log list

```hcl
data "tencentcloud_waf_attack_log_list" "example" {
  domain       = "domain.com"
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-07 00:00:00"
  query_string = "method:GET"
  sort         = "desc"
  query_count  = 10
  page         = 0
}
```

Obtain all domain name attack log list

```hcl
data "tencentcloud_waf_attack_log_list" "example" {
  domain       = "all"
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-07 00:00:00"
  query_string = "method:GET"
  sort         = "asc"
  query_count  = 20
  page         = 1
}
```