---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_deploy_config_version_history"
sidebar_current: "docs-tencentcloud-datasource-teo_deploy_config_version_history"
description: |-
  Use this data source to query detailed information of teo deploy config version history
---

# tencentcloud_teo_deploy_config_version_history

Use this data source to query detailed information of teo deploy config version history

## Example Usage

```hcl
data "tencentcloud_teo_deploy_config_version_history" "teo_deploy_config_version_history" {
  zone_id = "zone-2qtuhspy7cr6"
  env_id  = "env-2quhspyeq8r6"
}
```

## Argument Reference

The following arguments are supported:

* `env_id` - (Required, String) Environment ID.
* `zone_id` - (Required, String) Zone ID.
* `filters` - (Optional, List) Filtering condition. The maximum value of Filters.Values is 20. Detailed filtering conditions: record-id (Filter by release record ID).
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Value of the filtered field.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `records` - Release record details.


