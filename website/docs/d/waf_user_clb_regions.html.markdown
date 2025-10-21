---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_user_clb_regions"
sidebar_current: "docs-tencentcloud-datasource-waf_user_clb_regions"
description: |-
  Use this data source to query detailed information of waf user_clb_regions
---

# tencentcloud_waf_user_clb_regions

Use this data source to query detailed information of waf user_clb_regions

## Example Usage

```hcl
data "tencentcloud_waf_user_clb_regions" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Region list(ap-xxx format).
* `rich_datas` - Detail info for region.
  * `code` - Region code.
  * `id` - Region ID.
  * `text` - Chinese description for region.
  * `value` - English description for region.


