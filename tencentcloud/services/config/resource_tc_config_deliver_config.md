Provides a resource to manage Config delivery settings (global singleton configuration).

Example Usage

```hcl
resource "tencentcloud_config_deliver_config" "example" {
  status               = 1
  deliver_name         = "tf-example-deliver"
  target_arn           = "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"
  deliver_prefix       = "config"
  deliver_type         = "COS"
  deliver_content_type = 3
}
```
