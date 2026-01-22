---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_addons"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_addons"
description: |-
  Use this data source to query detailed information of kubernetes addons.
---

# tencentcloud_kubernetes_addons

Use this data source to query detailed information of kubernetes addons.

## Example Usage

```hcl
data "tencentcloud_kubernetes_addons" "kubernetes_addons" {
  cluster_id = "cls-12345678"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `addon_name` - (Optional, String) Add-on name (all add-ons in the cluster are returned if this parameter is not specified).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `addons` - List of add-ons.
  * `addon_name` - Add-on name.
  * `addon_version` - Add-on version.
  * `phase` - Add-on status
Note: This field may return `null`, indicating that no valid values can be obtained.
  * `raw_values` - Add-on parameters, which are base64-encoded strings in JSON/
Note: This field may return `null`, indicating that no valid values can be obtained.
  * `reason` - Reason for add-on failure
Note: This field may return `null`, indicating that no valid values can be obtained.


