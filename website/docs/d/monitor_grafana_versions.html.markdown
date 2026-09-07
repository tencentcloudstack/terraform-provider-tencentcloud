---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_versions"
sidebar_current: "docs-tencentcloud-datasource-monitor_grafana_versions"
description: |-
  Use this data source to query available Grafana versions of monitor
---

# tencentcloud_monitor_grafana_versions

Use this data source to query available Grafana versions of monitor

## Example Usage

```hcl
data "tencentcloud_monitor_grafana_versions" "grafana_versions" {
}

output "available_versions" {
  value = data.tencentcloud_monitor_grafana_versions.grafana_versions.versions
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `versions` - Grafana available version list.
  * `alias` - Grafana version alias.
  * `version` - Grafana version number.


