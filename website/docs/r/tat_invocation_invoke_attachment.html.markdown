---
subcategory: "TencentCloud Automation Tools(TAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tat_invocation_invoke_attachment"
sidebar_current: "docs-tencentcloud-resource-tat_invocation_invoke_attachment"
description: |-
  Provides a resource to create a tat invocation_invoke_attachment
---

# tencentcloud_tat_invocation_invoke_attachment

Provides a resource to create a tat invocation_invoke_attachment

## Example Usage

```hcl
resource "tencentcloud_tat_invocation_invoke_attachment" "invocation_invoke_attachment" {
  instance_id       = "ins-881b1c8w"
  working_directory = "/root"
  timeout           = 100
  # parameters = "{\"varA\": \"222\"}"
  username              = "root"
  output_cos_bucket_url = "https://BucketName-123454321.cos.ap-beijing.myqcloud.com"
  output_cos_key_prefix = "log"
  command_id            = "cmd-rxbs7f5z"
}
```

## Argument Reference

The following arguments are supported:

* `command_id` - (Required, String, ForceNew) Command ID.
* `instance_id` - (Required, String, ForceNew) ID of instances about to execute commands. Supported instance types:  CVM  LIGHTHOUSE.
* `output_cos_bucket_url` - (Optional, String, ForceNew) The COS bucket URL for uploading logs. The URL must start with https, such as https://BucketName-123454321.cos.ap-beijing.myqcloud.com.
* `output_cos_key_prefix` - (Optional, String, ForceNew) The COS bucket directory where the logs are saved; Check below for the rules of the directory name: 1 It must be a combination of number, letters, and visible characters, Up to 60 characters are allowed; 2 Use a slash (/) to create a subdirectory; 3 can not be used as the folder name; It cannot start with a slash (/), and cannot contain consecutive slashes.
* `parameters` - (Optional, String, ForceNew) Custom parameters of Command. The field type is JSON encoded string. For example, {varA: 222}.key is the name of the custom parameter and value is the default value. Both key and value are strings.If no parameter value is provided, the DefaultParameters is used.Up to 20 custom parameters are supported.The name of the custom parameter cannot exceed 64 characters and can contain [a-z], [A-Z], [0-9] and [-_].
* `timeout` - (Optional, Int, ForceNew) Command timeout period. Default value: 60 seconds. Value range: [1, 86400].
* `username` - (Optional, String, ForceNew) The username used to execute the command on the CVM or Lighthouse instance.The principle of least privilege is the best practice for permission management. We recommend you execute TAT commands as a general user. By default, the user root is used to execute commands on Linux and the user System is used on Windows.
* `working_directory` - (Optional, String, ForceNew) Command execution path. The default value is /root for SHELL commands and C:Program Filesqcloudtat_agentworkdir for POWERSHELL commands.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tat invocation can be imported using the invocation_id#instance_id, e.g.

```
terraform import tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment inv-mhs6ca8z#ins-881b1c8w
```

