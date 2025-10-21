---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_system_resource"
sidebar_current: "docs-tencentcloud-datasource-oceanus_system_resource"
description: |-
  Use this data source to query detailed information of oceanus system_resource
---

# tencentcloud_oceanus_system_resource

Use this data source to query detailed information of oceanus system_resource

## Example Usage

```hcl
data "tencentcloud_oceanus_system_resource" "example" {
  resource_ids = ["resource-abd503yt"]
  filters {
    name   = "Name"
    values = ["tf_example"]
  }
  cluster_id    = "cluster-n8yaia0p"
  flink_version = "Flink-1.11"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional, String) Cluster ID.
* `filters` - (Optional, List) Query the resource configuration list. If not specified, return all job configuration lists under ResourceIds.N.
* `flink_version` - (Optional, String) Query built-in connectors for the corresponding Flink version.
* `resource_ids` - (Optional, Set: [`String`]) Array of resource IDs to be queried.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Filter values for the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `resource_set` - Collection of resource details.
  * `latest_resource_config_version` - Latest version of the resource.
  * `name` - Resource name.
  * `region` - Region to which the resource belongs.
  * `remark` - Resource remarks.
  * `resource_id` - Resource ID.
  * `resource_type` - Resource type. 1 indicates JAR package, which is currently the only supported value.


