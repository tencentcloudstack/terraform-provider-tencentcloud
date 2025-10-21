---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_email_identities"
sidebar_current: "docs-tencentcloud-datasource-ses_email_identities"
description: |-
  Use this data source to query detailed information of ses email_identities
---

# tencentcloud_ses_email_identities

Use this data source to query detailed information of ses email_identities

## Example Usage

```hcl
data "tencentcloud_ses_email_identities" "email_identities" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `email_identities` - Sending domain name list.
  * `current_reputation_level` - Current credit rating.
  * `daily_quota` - Highest number of letters of the day.
  * `identity_name` - Sending domain name.
  * `identity_type` - Authentication type, fixed as DOMAIN.
  * `sending_enabled` - Is it verified.
* `max_daily_quota` - Maximum daily sending volume for a single domain name.
* `max_reputation_level` - Maximum credit rating.


