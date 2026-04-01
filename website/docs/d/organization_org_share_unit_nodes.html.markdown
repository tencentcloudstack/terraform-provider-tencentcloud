---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_share_unit_nodes"
sidebar_current: "docs-tencentcloud-datasource-organization_org_share_unit_nodes"
description: |-
  Use this data source to query organization org share unit nodes
---

# tencentcloud_organization_org_share_unit_nodes

Use this data source to query organization org share unit nodes

## Example Usage

```hcl
data "tencentcloud_organization_org_share_unit_nodes" "example" {
  unit_id = "us-xxxxx"
}
```

### Example with search_key:

```hcl
data "tencentcloud_organization_org_share_unit_nodes" "example" {
  unit_id    = "us-xxxxx"
  search_key = "123456"
}
```

## Argument Reference

The following arguments are supported:

* `unit_id` - (Required, String) Shared unit ID.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Search key, supports searching by department ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - List of share unit nodes.
  * `create_time` - Create time.
  * `share_node_id` - Department ID.


