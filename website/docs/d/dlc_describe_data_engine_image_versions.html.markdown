---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_data_engine_image_versions"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_data_engine_image_versions"
description: |-
  Use this data source to query detailed information of DLC describe data engine image versions
---

# tencentcloud_dlc_describe_data_engine_image_versions

Use this data source to query detailed information of DLC describe data engine image versions

## Example Usage

```hcl
data "tencentcloud_dlc_describe_data_engine_image_versions" "example" {
  engine_type = "SparkBatch"
  sort        = "UpdateTime"
  asc         = false
}
```

## Argument Reference

The following arguments are supported:

* `engine_type` - (Required, String) Engine type only support: SparkSQL/PrestoSQL/SparkBatch.
* `asc` - (Optional, Bool) Sort by: false (descending, default), true (ascending).
* `result_output_file` - (Optional, String) Used to save results.
* `sort` - (Optional, String) Sort fields: InsertTime (insert time, default), UpdateTime (update time).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `image_parent_versions` - Major version of the image information list of clusters.
  * `description` - Description of the major version of the image.
  * `engine_type` - Cluster types: SparkSQL, PrestoSQL, and SparkBatch.
  * `image_version_id` - ID of the major version of the image.
  * `image_version` - Name of the major version of the image.
  * `insert_time` - Insert time.
  * `is_public` - Whether it is a public version: 1: public version; 2: private version.
  * `is_shared_engine` - Version status. 1: initializing; 2: online; 3: offline.
  * `state` - Version status. 1: initializing; 2: online; 3: offline.
  * `update_time` - Update time.


