Provides a resource to create a tat invocation_invoke_attachment

Example Usage

```hcl
resource "tencentcloud_tat_invocation_invoke_attachment" "invocation_invoke_attachment" {
  instance_id = "ins-881b1c8w"
  working_directory = "/root"
  timeout = 100
  # parameters = "{\"varA\": \"222\"}"
  username = "root"
  output_cos_bucket_url = "https://BucketName-123454321.cos.ap-beijing.myqcloud.com"
  output_cos_key_prefix = "log"
  command_id = "cmd-rxbs7f5z"
}
```

Import

tat invocation can be imported using the invocation_id#instance_id, e.g.

```
terraform import tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment inv-mhs6ca8z#ins-881b1c8w
```