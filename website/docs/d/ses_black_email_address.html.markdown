---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_black_email_address"
sidebar_current: "docs-tencentcloud-datasource-ses_black_email_address"
description: |-
  Use this data source to query detailed information of ses black_email_address
---

# tencentcloud_ses_black_email_address

Use this data source to query detailed information of ses black_email_address

## Example Usage

```hcl
data "tencentcloud_ses_black_email_address" "black_email_address" {
  start_date    = "2020-09-22"
  end_date      = "2020-09-23"
  email_address = "xxx@mail.qcloud.com"
  task_id       = "7000"
}
```

## Argument Reference

The following arguments are supported:

* `end_date` - (Required, String) End date in the format of `YYYY-MM-DD`.
* `start_date` - (Required, String) Start date in the format of `YYYY-MM-DD`.
* `email_address` - (Optional, String) You can specify an email address to query.
* `result_output_file` - (Optional, String) Used to save results.
* `task_id` - (Optional, String) You can specify a task ID to query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `black_list` - List of blocklisted addresses.
  * `bounce_time` - Time when the email address is blocklisted.
  * `email_address` - Blocklisted email address.


