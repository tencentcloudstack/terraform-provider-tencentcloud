---
subcategory: "TencentCloud Automation Tools(TAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tat_command"
sidebar_current: "docs-tencentcloud-datasource-tat_command"
description: |-
  Use this data source to query detailed information of tat command
---

# tencentcloud_tat_command

Use this data source to query detailed information of tat command

## Example Usage

```hcl
data "tencentcloud_tat_command" "command" {
  # command_id = ""
  # command_name = ""
  command_type = "SHELL"
  created_by   = "TAT"
}
```

## Argument Reference

The following arguments are supported:

* `command_id` - (Optional, String) Command ID.
* `command_name` - (Optional, String) Command name.
* `command_type` - (Optional, String) Command type, Value is `SHELL` or `POWERSHELL`.
* `created_by` - (Optional, String) Command creator. `TAT` indicates a public command and `USER` indicates a personal command.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `command_set` - List of command details.
  * `command_id` - Command ID.
  * `command_name` - Command name.
  * `command_type` - Command type.
  * `content` - command.
  * `created_by` - Command creator. `TAT` indicates a public command and `USER` indicates a personal command.
  * `created_time` - Command creation time.
  * `default_parameters` - Default custom parameter value.
  * `description` - Command description.
  * `enable_parameter` - Whether to enable the custom parameter feature.
  * `formatted_description` - Formatted description of the command. This parameter is an empty string for user commands and contains values for public commands.
  * `output_cos_bucket_url` - The COS bucket URL for uploading logs.
  * `output_cos_key_prefix` - The COS bucket directory where the logs are saved.
  * `tags` - Tags bound to the command. At most 10 tags are allowed.
    * `key` - Tag key.
    * `value` - Tag value.
  * `timeout` - Command timeout period.
  * `updated_time` - Command update time.
  * `username` - The user who executes the command on the instance.
  * `working_directory` - Command execution path.


