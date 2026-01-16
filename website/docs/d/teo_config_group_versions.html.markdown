---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_config_group_versions"
sidebar_current: "docs-tencentcloud-datasource-teo_config_group_versions"
description: |-
  Use this data source to query detailed information of teo config group versions
---

# tencentcloud_teo_config_group_versions

Use this data source to query detailed information of teo config group versions

## Example Usage

```hcl
data "tencentcloud_teo_config_group_versions" "teo_config_group_versions" {
  zone_id  = "zone-2qtuhspy7cr6"
  group_id = "group-2quhspyeq8r6"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) Configuraration group ID.
* `zone_id` - (Required, String) Zone ID.
* `filters` - (Optional, List) Filtering condition. The maximum value of Filters.Values is 20. If this parameter is not specified, all version information for the selected configuration group is returned. Detailed filtering conditions: version-id (Filter by version ID).
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Value of the filtered field.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `config_group_version_infos` - Version information list.


