Provides a resource to create a tat invocation_command_attachment

Example Usage

```hcl
resource "tencentcloud_tat_invocation_command_attachment" "invocation_command_attachment" {
  content = base64encode("pwd")
  instance_id = "ins-881b1c8w"
  command_name = "terraform-test"
  description = "shell test"
  command_type = "SHELL"
  working_directory = "/root"
  timeout = 100
  save_command = false
  enable_parameter = false
  # default_parameters = "{\"varA\": \"222\"}"
  # parameters = "{\"varA\": \"222\"}"
  username = "root"
  output_cos_bucket_url = "https://BucketName-123454321.cos.ap-beijing.myqcloud.com"
  output_cos_key_prefix = "log"
}
```