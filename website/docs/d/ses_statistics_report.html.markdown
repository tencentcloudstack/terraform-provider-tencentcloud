---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_statistics_report"
sidebar_current: "docs-tencentcloud-datasource-ses_statistics_report"
description: |-
  Use this data source to query detailed information of ses statistics_report
---

# tencentcloud_ses_statistics_report

Use this data source to query detailed information of ses statistics_report

## Example Usage

```hcl
data "tencentcloud_ses_statistics_report" "statistics_report" {
  start_date             = "2020-10-01"
  end_date               = "2023-09-05"
  domain                 = "iac-tf.cloud"
  receiving_mailbox_type = "gmail.com"
}
```

## Argument Reference

The following arguments are supported:

* `end_date` - (Required, String) End date.
* `start_date` - (Required, String) Start date.
* `domain` - (Optional, String) Sender domain.
* `receiving_mailbox_type` - (Optional, String) Recipient address type, for example, gmail.com.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `daily_volumes` - Daily email sending statistics.
  * `accepted_count` - Number of email requests accepted by Tencent Cloud.
  * `bounce_count` - Number of bounced emails.
  * `clicked_count` - Number of recipients who clicked on links in emails.
  * `delivered_count` - Number of delivered emails.
  * `opened_count` - Number of users (deduplicated) who opened emails.
  * `request_count` - Number of email requests.
  * `send_date` - Date Note: this field may return null, indicating that no valid values can be obtained.
  * `unsubscribe_count` - Number of users who canceled subscriptions. Note: this field may return null, indicating that no valid values can be obtained.
* `overall_volume` - Overall email sending statistics.
  * `accepted_count` - Number of email requests accepted by Tencent Cloud.
  * `bounce_count` - Number of bounced emails.
  * `clicked_count` - Number of recipients who clicked on links in emails.
  * `delivered_count` - Number of delivered emails.
  * `opened_count` - Number of users (deduplicated) who opened emails.
  * `request_count` - Number of email requests.
  * `send_date` - Date Note: this field may return null, indicating that no valid values can be obtained.
  * `unsubscribe_count` - Number of users who canceled subscriptions. Note: this field may return null, indicating that no valid values can be obtained.


