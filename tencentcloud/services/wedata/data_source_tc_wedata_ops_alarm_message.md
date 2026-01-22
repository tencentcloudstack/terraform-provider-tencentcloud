Use this data source to query detailed information of wedata ops alarm message

Example Usage

```hcl
data "tencentcloud_wedata_ops_alarm_message" "wedata_ops_alarm_message" {
  project_id  = "1859317240494305280"
  alarm_message_id  = 263840
}
```