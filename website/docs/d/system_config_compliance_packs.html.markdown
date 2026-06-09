---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_system_config_compliance_packs"
sidebar_current: "docs-tencentcloud-datasource-system_config_compliance_packs"
description: |-
  Use this data source to query detailed information of Config system compliance packs.
---

# tencentcloud_system_config_compliance_packs

Use this data source to query detailed information of Config system compliance packs.

## Example Usage

### Query all system compliance packs

```hcl
data "tencentcloud_system_config_compliance_packs" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `compliance_pack_list` - System compliance pack list.
  * `compliance_pack_id` - Compliance pack ID.
  * `compliance_pack_name` - Compliance pack name.
  * `config_rules` - Config rules in the compliance pack.
    * `create_time` - Rule creation time.
    * `description` - Rule description.
    * `identifier` - Rule unique identifier.
    * `risk_level` - Rule risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
    * `rule_name` - Rule name.
    * `update_time` - Rule last update time.
  * `description` - Compliance pack description.
  * `risk_level` - Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).


