Provides a resource to create a tat invocation invoke attachment

Example Usage

```hcl
resource "tencentcloud_tat_invocation_invoke_attachment" "example" {
  instance_id       = "ins-hoek7x44"
  working_directory = "/root/"
  timeout           = 60
  username          = "root"
  command_id        = "cmd-l7otm4cn"
}
```

or

```hcl
resource "tencentcloud_tat_invocation_invoke_attachment" "example" {
  instance_id           = "ins-hoek7x44"
  working_directory     = "/root/"
  timeout               = 60
  username              = "root"
  command_id            = "cmd-l7otm4cn"
  output_cos_bucket_url = "https://your-bucket.cos.ap-guangzhou.myqcloud.com"
  output_cos_key_prefix = "tat/invoke"
}
```

Import

tat invocation can be imported using the invocation_id#instance_id, e.g.

```
terraform import tencentcloud_tat_invocation_invoke_attachment.example inv-64mrb10i1j#ins-hoek7x44
```