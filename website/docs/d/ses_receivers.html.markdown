---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_receivers"
sidebar_current: "docs-tencentcloud-datasource-ses_receivers"
description: |-
  Use this data source to query detailed information of ses receivers
---

# tencentcloud_ses_receivers

Use this data source to query detailed information of ses receivers

## Example Usage

```hcl
data "tencentcloud_ses_receivers" "receivers" {
  status   = 3
  key_word = "keep"
}
```

## Argument Reference

The following arguments are supported:

* `key_word` - (Optional, String) Group name keyword for fuzzy query.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, Int) Group status (`1`: to be uploaded; `2`: uploading; `3`: uploaded). To query groups in all states, do not pass in this parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Data record.
  * `count` - Total number of recipient email addresses.
  * `create_time` - Creation time, such as 2021-09-28 16:40:35.
  * `desc` - Recipient group descriptionNote: This field may return `null`, indicating that no valid value can be found.
  * `receiver_id` - Recipient group ID.
  * `receivers_name` - Recipient group name.
  * `receivers_status` - Group status (`1`: to be uploaded; `2` uploading; `3` uploaded)Note: This field may return `null`, indicating that no valid value can be found.


