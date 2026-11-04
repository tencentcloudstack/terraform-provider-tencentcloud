---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_list_aggregate_compliance_packs"
sidebar_current: "docs-tencentcloud-data_source-config_list_aggregate_compliance_packs"
description: |-
  Use this data source to query the list of Config aggregate compliance packs by account group.
---

# tencentcloud_config_list_aggregate_compliance_packs

Use this data source to query the list of Config aggregate compliance packs by account group.

## Example Usage

### Query all aggregate compliance packs of an account group

```hcl
data "tencentcloud_config_list_aggregate_compliance_packs" "example" {
  account_group_id = "ca-aggregator-abcdefg"
}
```

### Query aggregate compliance packs with filters

```hcl
data "tencentcloud_config_list_aggregate_compliance_packs" "example" {
  account_group_id     = "ca-aggregator-abcdefg"
  compliance_pack_name = "tf-example"
  risk_level           = [1, 2]
  status               = "ACTIVE"
  compliance_result    = ["NON_COMPLIANT"]
  order_type           = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `account_group_id` - (Required, String) Account group ID.
* `compliance_pack_name` - (Optional, String) Compliance pack name.
* `compliance_result` - (Optional, List) Compliance result. Valid values: COMPLIANT, NON_COMPLIANT.
* `include_compliance_pack_rule_result` - (Optional, String) Whether to include compliance pack rule result. Valid values: NO (not included). If not set or set to YES, the compliance result data is returned by default.
* `order_type` - (Optional, String) Sort type. Valid values: desc (descending), asc (ascending).
* `risk_level` - (Optional, List) Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
* `status` - (Optional, String) Compliance pack status. Valid values: ACTIVE, NO_ACTIVE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `compliance_pack_id` - Compliance pack ID.
  * `compliance_pack_name` - Compliance pack name.
  * `compliance_result` - Compliance result.
  * `create_time` - Creation time.
  * `description` - Compliance pack description.
  * `no_compliant_names` - List of non-compliant rule names.
  * `risk_level` - Risk level.
  * `rule_count` - Number of rules in the compliance pack.
  * `status` - Compliance pack status.
* `items` - Compliance pack list.
* `total` - Total number of compliance packs.


