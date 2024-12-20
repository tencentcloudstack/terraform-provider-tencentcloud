---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_share_units"
sidebar_current: "docs-tencentcloud-datasource-organization_org_share_units"
description: |-
  Use this data source to query detailed information of organization organization_org_share_units
---

# tencentcloud_organization_org_share_units

Use this data source to query detailed information of organization organization_org_share_units

## Example Usage

```hcl
data "tencentcloud_organization_org_share_units" "organization_org_share_units" {
  area       = "ap-guangzhou"
  search_key = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String) Shared unit area.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Search for keywords. Support UnitId and Name searches.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Shared unit list.


