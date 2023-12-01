---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_plugin_overviews"
sidebar_current: "docs-tencentcloud-datasource-monitor_grafana_plugin_overviews"
description: |-
  Use this data source to query detailed information of monitor grafana_plugin_overviews
---

# tencentcloud_monitor_grafana_plugin_overviews

Use this data source to query detailed information of monitor grafana_plugin_overviews

## Example Usage

```hcl
data "tencentcloud_monitor_grafana_plugin_overviews" "grafana_plugin_overviews" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `plugin_set` - Plugin set.
  * `plugin_id` - Grafana plugin ID.
  * `version` - Grafana plugin version.


