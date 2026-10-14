---
subcategory: "CDWPG"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwpg_nodes"
sidebar_current: "docs-tencentcloud-datasource-cdwpg_nodes"
description: |-
  Use this data source to query detailed information of cdwpg cdwpg_nodes
---

# tencentcloud_cdwpg_nodes

Use this data source to query detailed information of cdwpg cdwpg_nodes

## Example Usage

```hcl
data "tencentcloud_cdwpg_nodes" "cdwpg_nodes" {
  instance_id = "cdwpg-gexy9tue"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_nodes` - Node list.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeInstanceNodes](https://cloud.tencent.com/document/api/878/116657) - This API is used to list nodes.

