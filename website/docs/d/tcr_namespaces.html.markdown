---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_namespaces"
sidebar_current: "docs-tencentcloud-datasource-tcr_namespaces"
description: |-
  Use this data source to query detailed information of TCR namespaces.
---

# tencentcloud_tcr_namespaces

Use this data source to query detailed information of TCR namespaces.

## Example Usage

```hcl
data "tencentcloud_tcr_namespaces" "name" {
  instance_id    = "cls-satg5125"
  namespace_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the instance that the namespace belongs to.
* `namespace_name` - (Optional) ID of the TCR namespace to query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `namespace_list` - Information list of the dedicated TCR namespaces.
  * `is_public` - Indicate that the namespace is public or not.
  * `name` - Name of TCR namespace.


