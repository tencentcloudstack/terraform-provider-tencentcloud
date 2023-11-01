---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_instance_plugin_list"
sidebar_current: "docs-tencentcloud-datasource-elasticsearch_instance_plugin_list"
description: |-
  Use this data source to query detailed information of elasticsearch instance plugin list
---

# tencentcloud_elasticsearch_instance_plugin_list

Use this data source to query detailed information of elasticsearch instance plugin list

## Example Usage

```hcl
data "tencentcloud_elasticsearch_instance_plugin_list" "instance_plugin_list" {
  instance_id = "es-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `order_by_type` - (Optional, String) Order type. Valid values:
- asc: Ascending asc
- desc: Descending Desc.
* `order_by` - (Optional, String) order field. Valid values: `pluginName`.
* `plugin_type` - (Optional, Int) Plugin type. Valid values: `0`: System plugin.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `plugin_list` - Plugin information list.
  * `plugin_desc` - Plugin description.
  * `plugin_name` - Plugin name.
  * `plugin_type` - Plugin type. Valid values: `0`: System plugin.
  * `plugin_update_time` - Plugin update time.
  * `plugin_version` - Plugin version.
  * `removable` - Whether the plug-in can be uninstalled.
  * `status` - Plugin status. Valid values:
- `-2` has been uninstalled
- `-1` has been installed in
- `0` installation.


