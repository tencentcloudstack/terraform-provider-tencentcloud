---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_tag_retention_executions"
sidebar_current: "docs-tencentcloud-datasource-tcr_tag_retention_executions"
description: |-
  Use this data source to query detailed information of tcr tag_retention_executions
---

# tencentcloud_tcr_tag_retention_executions

Use this data source to query detailed information of tcr tag_retention_executions

## Example Usage

```hcl
data "tencentcloud_tcr_tag_retention_executions" "tag_retention_executions" {
  registry_id  = "tcr_ins_id"
  retention_id = 1
}
```

## Argument Reference

The following arguments are supported:

* `registry_id` - (Required, String) instance id.
* `retention_id` - (Required, Int) retention id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `retention_execution_list` - list of version retention execution records.
  * `end_time` - execution end time.
  * `execution_id` - execution id.
  * `retention_id` - retention id.
  * `start_time` - execution start time.
  * `status` - execution status: Failed, Succeed, Stopped, InProgress.


