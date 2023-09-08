---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_send_tasks"
sidebar_current: "docs-tencentcloud-datasource-ses_send_tasks"
description: |-
  Use this data source to query detailed information of ses send_tasks
---

# tencentcloud_ses_send_tasks

Use this data source to query detailed information of ses send_tasks

## Example Usage

```hcl
data "tencentcloud_ses_send_tasks" "send_tasks" {
  status      = 10
  receiver_id = 1063742
  task_type   = 1
}
```

## Argument Reference

The following arguments are supported:

* `receiver_id` - (Optional, Int) Recipient group ID.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, Int) Task status. `1`: to start; `5`: sending; `6`: sending suspended today; `7`: sending error; `10`: sent. To query tasks in all states, do not pass in this parameter.
* `task_type` - (Optional, Int) Task type. `1`: immediate; `2`: scheduled; `3`: recurring. To query tasks of all types, do not pass in this parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Data record.
  * `cache_count` - Number of emails cached.
  * `create_time` - Task creation time.
  * `cycle_param` - Parameters of a recurring taskNote: This field may return `null`, indicating that no valid value can be found.
    * `begin_time` - Start time of the task.
    * `interval_time` - Task recurrence in hours.
    * `term_cycle` - Specifies whether to end the cycle. This parameter is used to update the task. Valid values: 0: No; 1: Yes.
  * `err_msg` - Task exception informationNote: This field may return `null`, indicating that no valid value can be found.
  * `from_email_address` - Sender address.
  * `receiver_id` - Recipient group ID.
  * `receivers_name` - Recipient group name.
  * `request_count` - Number of emails requested to be sent.
  * `send_count` - Number of emails sent.
  * `subject` - Email subject.
  * `task_id` - Task ID.
  * `task_status` - Task status. `1`: to start; `5`: sending; `6`: sending suspended today; `7`: sending error; `10`: sent.
  * `task_type` - Task type. `1`: immediate; `2`: scheduled; `3`: recurring.
  * `template` - Template and template dataNote: This field may return `null`, indicating that no valid value can be found.
    * `template_data` - Variable parameters in the template. Please use `json.dump` to format the JSON object into a string type. The object is a set of key-value pairs. Each key denotes a variable, which is represented by {{key}}. The key will be replaced with the corresponding value (represented by {{value}}) when sending the email.Note: The parameter value cannot be data of a complex type such as HTML.Example: {name:xxx,age:xx}.
    * `template_id` - Template ID. If you do not have any template, please create one.
  * `timed_param` - Parameters of a scheduled taskNote: This field may return `null`, indicating that no valid value can be found.
    * `begin_time` - Start time of a scheduled sending task.
  * `update_time` - Task update time.


