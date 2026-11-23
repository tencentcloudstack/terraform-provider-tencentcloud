---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_async_job"
sidebar_current: "docs-tencentcloud-datasource-wedata_ops_async_job"
description: |-
  Use this data source to query detailed information of wedata ops workflow
---

# tencentcloud_wedata_ops_async_job

Use this data source to query detailed information of wedata ops workflow

## Example Usage

```hcl
data "tencentcloud_wedata_ops_async_job" "wedata_ops_async_job" {
  project_id = "1859317240494305280"
  async_id   = "9ba294ff-46d9-4a77-ae4a-acd0b4bdca3c"
}
```

## Argument Reference

The following arguments are supported:

* `async_id` - (Required, String) Asynchronous operation ID.
* `project_id` - (Required, String) Project ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Asynchronous operation detail result.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [GetOpsAsyncJob](https://www.tencentcloud.com/document/api/1174/73832)

