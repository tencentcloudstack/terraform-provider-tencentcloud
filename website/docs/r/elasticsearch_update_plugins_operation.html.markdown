---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_update_plugins_operation"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_update_plugins_operation"
description: |-
  Provides a resource to update elasticsearch plugins
---

# tencentcloud_elasticsearch_update_plugins_operation

Provides a resource to update elasticsearch plugins

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_update_plugins_operation" "update_plugins_operation" {
  instance_id         = "es-xxxxxx"
  install_plugin_list = ["analysis-pinyin"]
  force_restart       = false
  force_update        = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `force_restart` - (Optional, Bool, ForceNew) Whether to force a restart. Default is false.
* `force_update` - (Optional, Bool, ForceNew) Whether to reinstall, default value false.
* `install_plugin_list` - (Optional, Set: [`String`], ForceNew) List of plugins that need to be installed.
* `plugin_type` - (Optional, Int, ForceNew) Plugin type. 0: system plugin.
* `remove_plugin_list` - (Optional, Set: [`String`], ForceNew) List of plugins that need to be uninstalled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



