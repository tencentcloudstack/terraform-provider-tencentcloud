Use this data source to query detailed information of waf peak_points

Example Usage

Basic Query

```hcl
data "tencentcloud_waf_peak_points" "example" {
  from_time    = "2023-09-01 00:00:00"
  to_time      = "2023-09-07 00:00:00"
}
```

Query by filter

```hcl
data "tencentcloud_waf_peak_points" "example" {
  from_time   = "2023-09-01 00:00:00"
  to_time     = "2023-09-07 00:00:00"
  domain      = "domain.com"
  edition     = "clb-waf"
  instance_id = "waf_2kxtlbky00b2v1fn"
  metric_name = "access"
}
```