---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_product_namespace"
sidebar_current: "docs-tencentcloud-datasource-monitor_product_namespace"
description: |-
  Use this data source to query product namespace in monitor)
---

# tencentcloud_monitor_product_namespace

Use this data source to query product namespace in monitor)

## Example Usage

```hcl
data "tencentcloud_monitor_product_namespace" "instances" {
  name = "Redis"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name for filter, eg:`Load Banlancer`.
* `result_output_file` - (Optional) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list product namespaces. Each element contains the following attributes:
  * `namespace` - Namespace of each cloud product in monitor system.
  * `product_chinese_name` - Chinese name of this product.
  * `product_name` - English name of this product.


