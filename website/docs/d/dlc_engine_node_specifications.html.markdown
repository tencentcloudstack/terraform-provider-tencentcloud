---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_engine_node_specifications"
sidebar_current: "docs-tencentcloud-datasource-dlc_engine_node_specifications"
description: |-
  Use this data source to query detailed information of DLC engine node specifications
---

# tencentcloud_dlc_engine_node_specifications

Use this data source to query detailed information of DLC engine node specifications

## Example Usage

```hcl
data "tencentcloud_dlc_engine_node_specifications" "example" {
  data_engine_name = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_name` - (Optional, String) Engine Name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `driver_spec` - Driver available specifications.
* `executor_spec` - Available executor specifications.


