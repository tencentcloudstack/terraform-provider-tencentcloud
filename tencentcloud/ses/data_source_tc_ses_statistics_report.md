Use this data source to query detailed information of ses statistics_report

Example Usage

```hcl
data "tencentcloud_ses_statistics_report" "statistics_report" {
  start_date = "2020-10-01"
  end_date = "2023-09-05"
  domain = "iac-tf.cloud"
  receiving_mailbox_type = "gmail.com"
}
```