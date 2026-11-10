---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_list_aggregate_compliance_packs"
sidebar_current: "docs-tencentcloud-data_source-config_list_aggregate_compliance_packs"
description: |-
  Use this data source to query detailed information of TEO config aggregate compliance packs by account group.
---

# tencentcloud_config_list_aggregate_compliance_packs

Use this data source to query detailed information of TEO config aggregate compliance packs by account group.

## Example Usage

```hcl
data "tencentcloud_config_list_aggregate_compliance_packs" "example" {
  account_group_id = "aggr-xxxxxxxx"
}

output "first_pack_name" {
  value = data.tencentcloud_config_list_aggregate_compliance_packs.example.items[0].compliance_pack_name
}
```

### Query with filters:

```hcl
data "tencentcloud_config_list_aggregate_compliance_packs" "filtered" {
  account_group_id                    = "aggr-xxxxxxxx"
  compliance_pack_name                = "CIS-CWP-Best-Practices"
  risk_level                          = [1, 2]
  status                              = "ACTIVE"
  compliance_result                   = ["NON_COMPLIANT"]
  order_type                          = "desc"
  include_compliance_pack_rule_result = "NO"
}
```

## Argument Reference

The following arguments are supported:

* `account_group_id` - (Required, String) Account group ID.
* `compliance_pack_name` - (Optional, String) Compliance pack name used for filtering.
* `compliance_result` - (Optional, Set) Evaluation results used for filtering. Valid values: COMPLIANT, NON_COMPLIANT.
* `include_compliance_pack_rule_result` - (Optional, String) Whether to include compliance pack rule result. Valid values: NO. Default: empty. Note that empty or YES both return the compliance result data.
* `order_type` - (Optional, String) Sort type. Valid values: desc, asc.
* `risk_level` - (Optional, Set) Risk levels used for filtering. Valid values: 1 (high), 2 (medium), 3 (low).
* `status` - (Optional, String) Compliance pack status used for filtering. Valid values: ACTIVE, NO_ACTIVE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `compliance_pack_id` - Compliance pack ID.
  * `compliance_pack_name` - Compliance pack name.
  * `compliance_result` - Evaluation result of the compliance pack.
  * `create_time` - Creation time of the compliance pack.
  * `description` - Compliance pack description.
  * `no_compliant_names` - Names of non-compliant rules.
  * `risk_level` - Risk level. Valid values: 1 (high), 2 (medium), 3 (low).
  * `rule_count` - Number of rules in the compliance pack.
  * `status` - Compliance pack status.
* `items` - List of compliance packs.
* `total` - Total count of compliance packs matched by the filter.


