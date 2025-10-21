---
subcategory: "TencentCloud Automation Tools(TAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tat_invocation_command_attachment"
sidebar_current: "docs-tencentcloud-resource-tat_invocation_command_attachment"
description: |-
  Provides a resource to create a tat invocation_command_attachment
---

# tencentcloud_tat_invocation_command_attachment

Provides a resource to create a tat invocation_command_attachment

## Example Usage

```hcl
resource "tencentcloud_tat_invocation_command_attachment" "invocation_command_attachment" {
  content           = base64encode("pwd")
  instance_id       = "ins-881b1c8w"
  command_name      = "terraform-test"
  description       = "shell test"
  command_type      = "SHELL"
  working_directory = "/root"
  timeout           = 100
  save_command      = false
  enable_parameter  = false
  # default_parameters = "{\"varA\": \"222\"}"
  # parameters = "{\"varA\": \"222\"}"
  username              = "root"
  output_cos_bucket_url = "https://BucketName-123454321.cos.ap-beijing.myqcloud.com"
  output_cos_key_prefix = "log"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String, ForceNew) Base64-encoded command. The maximum length is 64 KB.
* `instance_id` - (Required, String, ForceNew) ID of instances about to execute commands. Supported instance types:  CVM  LIGHTHOUSE.
* `command_name` - (Optional, String, ForceNew) Command name. The name can be up to 60 bytes, and contain [a-z], [A-Z], [0-9] and [_-.].
* `command_type` - (Optional, String, ForceNew) Command type. SHELL and POWERSHELL are supported. The default value is SHELL.
* `default_parameters` - (Optional, String, ForceNew) The default value of the custom parameter value when it is enabled. The field type is JSON encoded string. For example, {varA: 222}.key is the name of the custom parameter and value is the default value. Both key and value are strings.If Parameters is not provided, the default values specified here are used.Up to 20 custom parameters are supported.The name of the custom parameter cannot exceed 64 characters and can contain [a-z], [A-Z], [0-9] and [-_].
* `description` - (Optional, String, ForceNew) Command description. The maximum length is 120 characters.
* `enable_parameter` - (Optional, Bool, ForceNew) Whether to enable the custom parameter feature.This cannot be modified once created.Default value: false.
* `output_cos_bucket_url` - (Optional, String, ForceNew) The COS bucket URL for uploading logs; The URL must start with https, such as https://BucketName-123454321.cos.ap-beijing.myqcloud.com.
* `output_cos_key_prefix` - (Optional, String, ForceNew) The COS bucket directory where the logs are saved; Check below for the rules of the directory name: 1 It must be a combination of number, letters, and visible characters, Up to 60 characters are allowed; 2 Use a slash (/) to create a subdirectory; 3 can not be used as the folder name; It cannot start with a slash (/), and cannot contain consecutive slashes.
* `parameters` - (Optional, String, ForceNew) Custom parameters of Command. The field type is JSON encoded string. For example, {varA: 222}.key is the name of the custom parameter and value is the default value. Both key and value are strings.If no parameter value is provided, the DefaultParameters is used.Up to 20 custom parameters are supported.The name of the custom parameter cannot exceed 64 characters and can contain [a-z], [A-Z], [0-9] and [-_].
* `save_command` - (Optional, Bool, ForceNew) Whether to save the command. Valid values:rue: SaveFalse:Do not saveThe default value is False.
* `timeout` - (Optional, Int, ForceNew) Command timeout period. Default value: 60 seconds. Value range: [1, 86400].
* `username` - (Optional, String, ForceNew) The username used to execute the command on the CVM or Lighthouse instance.The principle of least privilege is the best practice for permission management. We recommend you execute TAT commands as a general user. By default, the user root is used to execute commands on Linux and the user System is used on Windows.
* `working_directory` - (Optional, String, ForceNew) Command execution path. The default value is /root for SHELL commands and C:Program Filesqcloudtat_agentworkdir for POWERSHELL commands.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `command_id` - Command ID.


