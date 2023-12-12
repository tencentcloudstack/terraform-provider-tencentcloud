Use this data source to query detailed information of css time_shift_stream_list

Example Usage

```hcl
data "tencentcloud_css_time_shift_stream_list" "time_shift_stream_list" {
  start_time   = 1698768000
  end_time     = 1698820641
  stream_name  = "live"
  domain       = "177154.push.tlivecloud.com"
  domain_group = "tf-test"
}
```