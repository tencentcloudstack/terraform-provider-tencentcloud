---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_products"
sidebar_current: "docs-tencentcloud-datasource-ssm_products"
description: |-
  Use this data source to query detailed information of ssm products
---

# tencentcloud_ssm_products

Use this data source to query detailed information of ssm products

## Example Usage

```hcl
data "tencentcloud_ssm_products" "products" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `products` - List of supported services.


