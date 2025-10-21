---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_domains"
sidebar_current: "docs-tencentcloud-datasource-css_domains"
description: |-
  Use this data source to query detailed information of css domains
---

# tencentcloud_css_domains

Use this data source to query detailed information of css domains

## Example Usage

```hcl
data "tencentcloud_css_domains" "domains" {
  domain_type   = 0
  play_type     = 1
  is_delay_live = 0
}
```

## Argument Reference

The following arguments are supported:

* `domain_prefix` - (Optional, String) domain name prefix.
* `domain_status` - (Optional, Int) domain name status filter. 0-disable, 1-enable.
* `domain_type` - (Optional, Int) Domain name type filtering. 0-push, 1-play.
* `is_delay_live` - (Optional, Int) 0 normal live broadcast 1 slow live broadcast default 0.
* `play_type` - (Optional, Int) Playing area, this parameter is meaningful only when DomainType=1. 1: Domestic.2: Global.3: Overseas.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domain_list` - A list of domain name details.
  * `b_c_name` - Is there a CName to the fixed rule domain name: 0: No. 1: Yes.
  * `create_time` - add time.Note: This field is Beijing time (UTC+8 time zone).
  * `current_c_name` - The cname information used by the current client.
  * `is_delay_live` - Whether to slow live broadcast: 0: normal live broadcast. 1: Slow live broadcast.
  * `is_mini_program_live` - 0: Standard live broadcast. 1: Mini program live broadcast. Note: This field may return null, indicating that no valid value can be obtained.
  * `name` - Live domain name.
  * `play_type` - Playing area, this parameter is meaningful only when Type=1. 1: Domestic. 2: Global. 3: Overseas.
  * `rent_expire_time` - Failure parameter, can be ignored. Note: This field is Beijing time (UTC+8 time zone).
  * `rent_tag` - invalid parameter, can be ignored.
  * `status` - Domain Status: 0: disable. 1: Enabled.
  * `target_domain` - The domain name corresponding to the cname.
  * `type` - Domain Type: 0: push stream. 1: Play.


