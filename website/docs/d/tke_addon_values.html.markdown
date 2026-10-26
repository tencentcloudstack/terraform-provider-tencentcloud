---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tke_addon_values"
sidebar_current: "docs-tencentcloud-datasource-tke_addon_values"
description: |-
  Use this data source to query detailed information of TKE addon values.
---

# tencentcloud_tke_addon_values

Use this data source to query detailed information of TKE addon values.

## Example Usage

```hcl
data "tencentcloud_tke_addon_values" "example" {
  cluster_id = "cls-5yezvaxo"
  addon_name = "ip-masq-agent"
}

locals {
  current_values   = jsondecode(data.tencentcloud_tke_addon_values.example.values)
  available_values = jsondecode(data.tencentcloud_tke_addon_values.example.default_values)
}

output "addon_current_values" {
  value = local.current_values
}

output "addon_available_values" {
  value = local.available_values
}
```

## Argument Reference

The following arguments are supported:

* `addon_name` - (Required, String) Name of the addon, which can be obtained from the component management page in the cluster details page of the TKE console.
* `cluster_id` - (Required, String) ID of the TKE cluster, which can be obtained from the cluster list in the TKE console (https://console.cloud.tencent.com/tke2/cluster).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `default_values` - Parameter list supported by the addon in JSON format, whose value is the default value of the chart.
* `values` - Parameter list in JSON format. If the addon is installed, the parameters are rendered with the configured parameters and the default parameters of the chart; otherwise the value is empty.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeAddonValues](https://www.tencentcloud.com/document/api/457/56952)

