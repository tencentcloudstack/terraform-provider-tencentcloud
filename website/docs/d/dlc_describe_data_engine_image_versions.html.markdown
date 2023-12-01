---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_data_engine_image_versions"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_data_engine_image_versions"
description: |-
  Use this data source to query detailed information of dlc describe_data_engine_image_versions
---

# tencentcloud_dlc_describe_data_engine_image_versions

Use this data source to query detailed information of dlc describe_data_engine_image_versions

## Example Usage

```hcl
data "tencentcloud_dlc_describe_data_engine_image_versions" "describe_data_engine_image_versions" {
  engine_type = "SparkBatch"
}
```

## Argument Reference

The following arguments are supported:

* `engine_type` - (Required, String) Engine type only support: SparkSQL/PrestoSQL/SparkBatch.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `image_parent_versions` - Cluster large version image information list.
  * `description` - Image major version description.
  * `engine_type` - Engine type only support: SparkSQL/PrestoSQL/SparkBatch.
  * `image_version_id` - Engine major version id.
  * `image_version` - Engine major version name.
  * `insert_time` - Create time.
  * `is_public` - Whether it is a public version, only support: 1: public;/2: private.
  * `is_shared_engine` - Is shared engine, only support: 1:yes/2:no.
  * `state` - Version status, only support: 1: initialized/2: online/3: offline.
  * `update_time` - Update time.


