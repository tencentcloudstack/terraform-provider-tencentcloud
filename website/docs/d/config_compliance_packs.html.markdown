---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_compliance_packs"
sidebar_current: "docs-tencentcloud-datasource-config_compliance_packs"
description: |-
  Use this data source to query detailed information of Config compliance packs.
---

# tencentcloud_config_compliance_packs

Use this data source to query detailed information of Config compliance packs.

## Example Usage

### Query all compliance packs

```hcl
data "tencentcloud_config_compliance_packs" "example" {}
```

### Query compliance packs by name

```hcl
data "tencentcloud_config_compliance_packs" "example" {
  compliance_pack_name = "tf-example"
}
```

### Query compliance packs by filters

```hcl
data "tencentcloud_config_compliance_packs" "example" {
  compliance_pack_name = "tf-example"
  risk_level           = [1, 2]
  status               = "ACTIVE"
  compliance_result    = ["NON_COMPLIANT"]
  order_type           = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `compliance_pack_name` - (Optional, String) Compliance pack name for filtering.
* `compliance_result` - (Optional, List: [`String`]) Compliance result list for filtering. Valid values: COMPLIANT, NON_COMPLIANT.
* `order_type` - (Optional, String) Sort type. Valid values: desc (descending), asc (ascending).
* `result_output_file` - (Optional, String) Used to save results.
* `risk_level` - (Optional, List: [`Int`]) Risk level list for filtering. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
* `status` - (Optional, String) Compliance pack status for filtering. Valid values: ACTIVE, NO_ACTIVE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `compliance_pack_list` - Compliance pack list.
  * `compliance_pack_id` - Compliance pack ID.
  * `compliance_pack_name` - Compliance pack name.
  * `compliance_result` - Compliance result. Valid values: COMPLIANT, NON_COMPLIANT.
  * `create_time` - Creation time.
  * `description` - Compliance pack description.
  * `no_compliant_names` - List of non-compliant rule names.
  * `risk_level` - Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
  * `rule_count` - Number of rules in the compliance pack.
  * `status` - Compliance pack status. Valid values: ACTIVE, NO_ACTIVE.


