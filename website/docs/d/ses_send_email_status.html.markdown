---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_send_email_status"
sidebar_current: "docs-tencentcloud-datasource-ses_send_email_status"
description: |-
  Use this data source to query detailed information of ses send_email_status
---

# tencentcloud_ses_send_email_status

Use this data source to query detailed information of ses send_email_status

## Example Usage

```hcl
data "tencentcloud_ses_send_email_status" "send_email_status" {
  request_date     = "2020-09-22"
  message_id       = "qcloudses-30-4123414323-date-20210101094334-syNARhMTbKI1"
  to_email_address = "example@cloud.com"
}
```

## Argument Reference

The following arguments are supported:

* `request_date` - (Required, String) Date sent. This parameter is required. You can only query the sending status for a single date at a time.
* `message_id` - (Optional, String) The MessageId field returned by the SendMail API.
* `result_output_file` - (Optional, String) Used to save results.
* `to_email_address` - (Optional, String) Recipient email address.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `email_status_list` - Status of sent emails.
  * `deliver_message` - Description of the recipient processing status.
  * `deliver_status` - Recipient processing status0: Tencent Cloud has accepted the request and added it to the send queue.1: The email is delivered successfully. DeliverTime indicates the time when the email is delivered successfully.2: The email is discarded. DeliverMessage indicates the reason for discarding.3: The recipient&amp;#39;s ESP rejects the email, probably because the email address does not exist or due to other reasons.8: The email is delayed by the ESP. DeliverMessage indicates the reason for delay.
  * `deliver_time` - Timestamp when Tencent Cloud delivers the email.
  * `from_email_address` - Sender email address.
  * `message_id` - The MessageId field returned by the SendEmail API.
  * `request_time` - Timestamp when the request arrives at Tencent Cloud.
  * `send_status` - Tencent Cloud processing status0: Successful.1001: Internal system exception.1002: Internal system exception.1003: Internal system exception.1003: Internal system exception.1004: Email sending timed out.1005: Internal system exception.1006: You have sent too many emails to the same address in a short period.1007: The email address is in the blocklist.1008: The sender domain is rejected by the recipient.1009: Internal system exception.1010: The daily email sending limit is exceeded.1011: You have no permission to send custom content. Use a template.1013: The sender domain is unsubscribed from by the recipient.2001: No results were found.3007: The template ID is invalid or the template is unavailable.3008: The sender domain is temporarily blocked by the recipient domain.3009: You have no permission to use this template.3010: The format of the TemplateData field is incorrect. 3014: The email cannot be sent because the sender domain is not verified.3020: The recipient email address is in the blocklist.3024: Failed to precheck the email address format.3030: Email sending is restricted temporarily due to a high bounce rate.3033: The account has insufficient balance or overdue payment.
  * `to_email_address` - Recipient email address.
  * `user_clicked` - Whether the recipient has clicked the links in the email.
  * `user_complainted` - Whether the recipient has reported the sender.
  * `user_opened` - Whether the recipient has opened the email.
  * `user_unsubscribed` - Whether the recipient has unsubscribed from the email sent by the sender.


