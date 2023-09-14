---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_find_domains"
sidebar_current: "docs-tencentcloud-datasource-waf_find_domains"
description: |-
  Use this data source to query detailed information of waf find_domains
---

# tencentcloud_waf_find_domains

Use this data source to query detailed information of waf find_domains

## Example Usage

### Find all domains

```hcl
data "tencentcloud_waf_find_domains" "example" {}
```

### Find domains by filter

```hcl
data "tencentcloud_waf_find_domains" "example" {
  key           = "keyWord"
  is_waf_domain = "1"
  by            = "FindTime"
  order         = "asc"
}
```

## Argument Reference

The following arguments are supported:

* `by` - (Optional, String) Sorting parameter, eg: FindTime.
* `is_waf_domain` - (Optional, String) Whether access to waf or not.
* `key` - (Optional, String) Filter condition.
* `order` - (Optional, String) Sorting type, eg: desc, asc.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Domain info list.
  * `appid` - User appid.
  * `domain_id` - Domain unique id.
  * `domain` - Domain name.
  * `edition` - Instance type, sparta-waf represents SAAS WAF, clb-waf represents CLB WAF.
  * `find_time` - Find time.
  * `instance_id` - Instance unique id.
  * `ips` - Domain ip.
  * `is_waf_domain` - Whether access to waf or not.


