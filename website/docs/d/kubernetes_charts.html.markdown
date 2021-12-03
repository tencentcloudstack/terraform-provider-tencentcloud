---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_charts"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_charts"
description: |-
  Use this data source to query detailed information of kubernetes cluster addons.
---

# tencentcloud_kubernetes_charts

Use this data source to query detailed information of kubernetes cluster addons.

## Example Usage

```hcl
data "tencentcloud_kubernetes_charts" "name" {}
```

## Argument Reference

The following arguments are supported:

* `arch` - (Optional) Operation system app supported. Available values: `arm32`, `arm64`, `amd64`.
* `cluster_type` - (Optional) Cluster type. Available values: `tke`, `eks`.
* `kind` - (Optional) Kind of app chart. Available values: `log`, `scheduler`, `network`, `storage`, `monitor`, `dns`, `image`, `other`, `invisible`.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `chart_list` - App chart list.
  * `label` - Label of chart.
  * `latest_version` - Chart latest version.
  * `name` - Name of chart.


