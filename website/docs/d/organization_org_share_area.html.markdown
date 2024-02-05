---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_share_area"
sidebar_current: "docs-tencentcloud-datasource-organization_org_share_area"
description: |-
  Use this data source to query detailed information of organization org_share_area
---

# tencentcloud_organization_org_share_area

Use this data source to query detailed information of organization org_share_area

## Example Usage

```hcl
data "tencentcloud_organization_org_share_area" "org_share_area" {
  lang = "zh"
}
```

## Argument Reference

The following arguments are supported:

* `lang` - (Optional, String) Language.default zh.
Valid values:
  - `zh`: Chinese.
  - `en`: English.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Area list.
  * `area_id` - Region ID.
  * `area` - Region identifier.
  * `name` - Region name.


