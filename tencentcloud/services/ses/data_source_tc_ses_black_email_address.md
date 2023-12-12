Use this data source to query detailed information of ses black_email_address

Example Usage

```hcl
data "tencentcloud_ses_black_email_address" "black_email_address" {
  start_date = "2020-09-22"
  end_date = "2020-09-23"
  email_address = "xxx@mail.qcloud.com"
  task_id = "7000"
}
```