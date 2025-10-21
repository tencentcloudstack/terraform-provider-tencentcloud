---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_share_unit_resources"
sidebar_current: "docs-tencentcloud-datasource-organization_org_share_unit_resources"
description: |-
  Use this data source to query detailed information of organization organization_org_share_unit_resources
---

# tencentcloud_organization_org_share_unit_resources

Use this data source to query detailed information of organization organization_org_share_unit_resources

## Example Usage

```hcl
data "tencentcloud_organization_org_share_unit_resources" "organization_org_share_unit_resources" {
  area    = "ap-guangzhou"
  unit_id = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String) Shared unit area.
* `unit_id` - (Required, String) Shared unit ID.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Search for keywords. Support product resource ID search.
* `type` - (Optional, String) Shared resource type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Shared unit resource list.


