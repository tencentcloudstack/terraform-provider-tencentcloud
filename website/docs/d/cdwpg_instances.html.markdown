---
subcategory: "CDWPG"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwpg_instances"
sidebar_current: "docs-tencentcloud-datasource-cdwpg_instances"
description: |-
  Use this data source to query detailed information of cdwpg cdwpg_instances
---

# tencentcloud_cdwpg_instances

Use this data source to query detailed information of cdwpg cdwpg_instances

## Example Usage

```hcl
data "tencentcloud_cdwpg_instances" "cdwpg_instances" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `search_instance_id` - (Optional, String) Search instance id.
* `search_instance_name` - (Optional, String) Search instance name.
* `search_tags` - (Optional, Set: [`String`]) Search tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances_list` - Instances list.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeSimpleInstances](https://cloud.tencent.com/document/api/878/100126) - This API is used to get a list of instance

