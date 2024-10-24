---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_route_table_input_policies"
sidebar_current: "docs-tencentcloud-datasource-ccn_route_table_input_policies"
description: |-
  Use this data source to query CCN route table input policies.
---

# tencentcloud_ccn_route_table_input_policies

Use this data source to query CCN route table input policies.

## Example Usage

```hcl
data "tencentcloud_ccn_route_table_input_policies" "example" {
  ccn_id         = "ccn-06jek8tf"
  route_table_id = "ccnrtb-4jv5ltb9"
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String) CCN Instance ID.
* `route_table_id` - (Required, String) CCN Route table ID.
* `policy_version` - (Optional, Int) Policy version.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policy_set` - Policy set.


