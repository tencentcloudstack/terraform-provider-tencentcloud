---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_job_status_detail"
sidebar_current: "docs-tencentcloud-datasource-emr_job_status_detail"
description: |-
  Use this data source to query detailed information of emr emr_job_status_detail
---

# tencentcloud_emr_job_status_detail

Use this data source to query detailed information of emr emr_job_status_detail

## Example Usage

```hcl
data tencentcloud_emr_job_status_detail "emr_job_status_detail" {
  instance_id = "emr-byhnjsb3"
  flow_param {
    f_key   = "FlowId"
    f_value = "1921228"
  }
}
```

## Argument Reference

The following arguments are supported:

* `flow_param` - (Required, List) Flow-related Parameters.
* `instance_id` - (Required, String) EMR Instance ID.
* `need_extra_detail` - (Optional, Bool) Whether to return additional task information.
* `result_output_file` - (Optional, String) Used to save results.

The `flow_param` object supports the following:

* `f_key` - (Required, String) Process Parameter Key: value range: TraceId: Query by TraceId FlowId: Query by FlowId.
* `f_value` - (Required, String) Parameter Value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `flow_desc` - Flow Parameter Description.
  * `p_key` - Parameter Key.
  * `p_value` - Parameter Value.
* `flow_extra_detail` - Flow Extra Execution Detail,Return when NeedExtraDetail is true.
  * `detail` - Flow Extra Execution Detail.
    * `p_key` - Parameter Key.
    * `p_value` - Parameter Value.
  * `title` - Flow Extra Execution Detail Title.
* `flow_name` - Flow Name.
* `flow_total_progress` - Flow Total Execution Progress.
* `flow_total_status` - Flow Total Execution Status, 0: Initialized, 1: Running, 2: Completed, 3: Completed (with skipped steps), -1: Failed, -3: Blocke.
* `stage_details` - Task Information.
  * `desc` - Flow Execution Status Description.
  * `endtime` - Flow Execution End Time.
  * `failed_reason` - Flow Execution Failure Reason.
  * `had_wood_detail` - Whether to return additional task information.
  * `is_show` - Whether to display the flow.
  * `is_sub_flow` - Whether it is a sub-flow.
  * `language_key` - Multilingual Version Key.
  * `name` - Step Name.
  * `progress` - Flow Execution Progress.
  * `stage` - Step ID.
  * `starttime` - Flow Execution Start Time.
  * `status` - Flow Execution Status: 0: Not Started, 1: In Progress, 2: Completed, 3: Partially Completed, -1: Failed.
  * `sub_flow_flag` - Sub-Flow Flag.
  * `time_consuming` - Flow Execution Time Consuming.
  * `wood_job_id` - Wood Subprocess ID.


