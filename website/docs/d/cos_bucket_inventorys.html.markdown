---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_bucket_inventorys"
sidebar_current: "docs-tencentcloud-datasource-cos_bucket_inventorys"
description: |-
  Use this data source to query the COS bucket inventorys.
---

# tencentcloud_cos_bucket_inventorys

Use this data source to query the COS bucket inventorys.

## Example Usage

```hcl
data "tencentcloud_cos_bucket_inventorys" "cos_bucket_inventorys" {
  bucket = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) Bucket.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `inventorys` - Multiple batch processing task information.
  * `destination` - Information about the inventory result destination.
    * `account_id` - ID of the bucket owner.
    * `bucket` - Bucket name.
    * `encryption` - Server-side encryption for the inventory result.
      * `sse_cos` - Encryption with COS-managed key. This field can be left empty.
    * `format` - Format of the inventory result. Valid value: CSV.
    * `prefix` - Prefix of the inventory result.
  * `filter` - Filters objects prefixed with the specified value to analyze.
    * `period` - Creation time range of the objects to analyze.
      * `end_time` - Creation end time of the objects to analyze. The parameter is a timestamp in seconds, for example, 1568688762.
      * `start_time` - Creation start time of the objects to analyze. The parameter is a timestamp in seconds, for example, 1568688761.
    * `prefix` - Prefix of the objects to analyze.
  * `id` - Whether to enable the inventory. true or false.
  * `included_object_versions` - Whether to include object versions in the inventory. All or No.
  * `is_enabled` - Whether to enable the inventory. true or false.
  * `schedule` - Inventory job cycle.
    * `frequency` - Frequency of the inventory job. Enumerated values: Daily, Weekly.


