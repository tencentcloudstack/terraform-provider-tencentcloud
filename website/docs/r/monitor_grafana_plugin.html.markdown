---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_plugin"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_plugin"
description: |-
  Provides a resource to create a monitor grafanaPlugin
---

# tencentcloud_monitor_grafana_plugin

Provides a resource to create a monitor grafanaPlugin

## Example Usage

```hcl
resource "tencentcloud_monitor_grafana_plugin" "grafanaPlugin" {
  instance_id = "grafana-50nj6v00"
  plugin_id   = "grafana-piechart-panel"
  version     = "1.6.2"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Grafana instance id.
* `plugin_id` - (Required, String) Plugin id.
* `version` - (Optional, String) Plugin version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor grafanaPlugin can be imported using the instance_id#plugin_id, e.g.
```
$ terraform import tencentcloud_monitor_grafana_plugin.grafanaPlugin grafana-50nj6v00#grafana-piechart-panel
```

