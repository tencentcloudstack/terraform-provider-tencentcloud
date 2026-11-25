---
subcategory: "CSIP"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_csip_check_view_risks"
sidebar_current: "docs-tencentcloud-data_source-csip_check_view_risks"
description: |-
  Use this data source to query the cloud asset configuration risks under the check item view.
---

# tencentcloud_csip_check_view_risks

Use this data source to query the cloud asset configuration risks under the check item view.

## Example Usage

```hcl
data "tencentcloud_csip_check_view_risks" "example" {
  member_id = ["mem-6wfo0fzks3"]
  filters {
    name   = "Status"
    values = ["CHECKING"]
  }
  order = "desc"
  by    = "update_time"
}
```

## Argument Reference

The following arguments are supported:

* `by` - (Optional, String) Sort field.
* `filters` - (Optional, List) Filter conditions.
* `member_id` - (Optional, List) Member IDs of the organization account.
* `order` - (Optional, String) Sort type.

The `filters` object supports the following:

* `name` - (Required, String) Filter name. Valid value: Status (execution result, Values: SUCCESS/FAILED/USER_CANCELED/CHECKING).
* `values` - (Required, List) List of filter values.
* `exact_match` - (Optional, String) Whether to match exactly: 1 means exact match; fuzzy match by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `tag` - Tag.
    * `terms` - Terms.
  * `asset_count` - Number of affected assets.
  * `asset_type_icon_url` - Asset type icon.
  * `asset_type` - Asset type.
  * `check_type` - Check type.
  * `classify` - Disposition classification.
  * `create_time` - First discovery time.
  * `event_type` - Event type.
  * `id` - Standard ID.
  * `name` - Standard name.
  * `provider` - Cloud provider.
  * `risk_count` - Number of risks.
  * `risk_desc` - 1 risk item exists.
  * `risk_rule_id` - Check item rule ID.
  * `risk_status` - Risk status.
  * `risk_title` - Risk name.
  * `severity` - Risk level.
  * `standard_terms` - CSPM standard terms.
  * `text` - CVM instance type.
  * `text` - CVM instance type.
  * `update_time` - Risk update time.
  * `value` - CVM instance name.
  * `value` - CVM instance name.
* `asset_type_list` - Collection of asset types.
* `check_view_risk_list` - List of risks under the check item view.
* `id` - ID of the data source, a random token.
* `provider_list` - Collection of cloud provider types.
* `standard_name_list` - List of CSPM standard tags under the check item view.


