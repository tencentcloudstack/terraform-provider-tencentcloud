---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_ext_namespace"
sidebar_current: "docs-tencentcloud-datasource-monitor_ext_namespace"
description: |-
  Use this data source to query Monitor ext namespaces.
---

# tencentcloud_monitor_ext_namespace

Use this data source to query Monitor ext namespaces.

~> **NOTE:** This data source queries all ext namespaces returned by the `DescribeExtNamespace` API. The API has no request parameters and no pagination, so the result set is not filterable.

## Example Usage

```hcl
data "tencentcloud_monitor_ext_namespace" "example" {}
```
