---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_share_unit_members"
sidebar_current: "docs-tencentcloud-datasource-organization_org_share_unit_members"
description: |-
  Use this data source to query detailed information of organization organization_org_share_unit_members
---

# tencentcloud_organization_org_share_unit_members

Use this data source to query detailed information of organization organization_org_share_unit_members

## Example Usage

```hcl
data "tencentcloud_organization_org_share_unit_members" "organization_org_share_unit_members" {
  unit_id    = "xxxxxx"
  area       = "ap-guangzhou"
  search_key = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String) Shared unit area.
* `unit_id` - (Required, String) Shared unit ID.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Search for keywords. Support member Uin searches.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Shared unit member list.


