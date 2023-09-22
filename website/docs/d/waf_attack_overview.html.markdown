---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_attack_overview"
sidebar_current: "docs-tencentcloud-datasource-waf_attack_overview"
description: |-
  Use this data source to query detailed information of waf attack_overview
---

# tencentcloud_waf_attack_overview

Use this data source to query detailed information of waf attack_overview

## Example Usage

### Basic Query

```hcl
data "tencentcloud_waf_attack_overview" "example" {
  from_time = "2023-09-01 00:00:00"
  to_time   = "2023-09-07 00:00:00"
}
```

### Query by filter

```hcl
data "tencentcloud_waf_attack_overview" "example" {
  from_time   = "2023-09-01 00:00:00"
  to_time     = "2023-09-07 00:00:00"
  appid       = 1304251372
  domain      = "test.com"
  edition     = "clb-waf"
  instance_id = "waf_2kxtlbky00b2v1fn"
}
```

## Argument Reference

The following arguments are supported:

* `from_time` - (Required, String) Begin time.
* `to_time` - (Required, String) End time.
* `appid` - (Optional, Int) App id.
* `domain` - (Optional, String) Domain.
* `edition` - (Optional, String) support `sparta-waf`, `clb-waf`, otherwise not filter.
* `instance_id` - (Optional, String) Waf instanceId, otherwise not filter.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_count` - Access count.
* `acl_count` - Access control count.
* `api_assets_count` - Api asset count.
* `api_risk_event_count` - Number of API risk events.
* `attack_count` - Attack count.
* `bot_count` - Bot attack count.
* `cc_count` - CC attack count.


