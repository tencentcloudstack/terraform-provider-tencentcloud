Use this data source to query detailed information of css time_shift_record_detail

Example Usage

```hcl
data "tencentcloud_css_time_shift_record_detail" "time_shift_record_detail" {
  domain        = "177154.push.tlivecloud.com"
  app_name      = "qqq"
  stream_name   = "live"
  start_time    = 1698768000
  end_time      = 1698820641
  domain_group  = "tf-test"
  trans_code_id = 0
}
```