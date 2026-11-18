subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_ext_metric"
sidebar_current: "docs-tencentcloud-datasource-monitor_ext_metric"
description: |-
  Use this data source to query Cloud Monitor external metrics of an external namespace.
---

# tencentcloud_monitor_ext_metric

Use this data source to query Cloud Monitor external metrics of an external namespace.

## Example Usage

```hcl
data "tencentcloud_monitor_ext_metric" "example" {
  ext_namespace = "qcs::monitor:my-namespace"
}
```
