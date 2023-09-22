---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_user_domains"
sidebar_current: "docs-tencentcloud-datasource-waf_user_domains"
description: |-
  Use this data source to query detailed information of waf user_domains
---

# tencentcloud_waf_user_domains

Use this data source to query detailed information of waf user_domains

## Example Usage

```hcl
data "tencentcloud_waf_user_domains" "user_domains" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `users_info` - Domain infos.
  * `appid` - User appid.
  * `cls` - CLS switch 1: write, 0: do not writeNote: This field may return null, indicating that a valid value cannot be obtained.
  * `domain_id` - Domain unique id.
  * `domain` - Domain name.
  * `edition` - Instance type, sparta-waf represents SAAS WAF, clb-waf represents CLB WAF.
  * `instance_id` - Instance unique id.
  * `instance_name` - Instance name.
  * `level` - Instance level infoNote: This field may return null, indicating that a valid value cannot be obtained.
  * `write_config` - Switch for accessing log fieldsNote: This field may return null, indicating that a valid value cannot be obtained.


