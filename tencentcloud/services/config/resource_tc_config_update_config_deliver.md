Provides a resource to manage Tencent Cloud Config delivery settings (global singleton configuration).

Example Usage

```hcl
resource "tencentcloud_config_update_config_deliver" "example" {
  status               = 1
  deliver_name         = "tf-example-deliver"
  target_arn           = "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"
  deliver_prefix       = "config"
  deliver_type         = "COS"
  deliver_content_type = 3
}
```

## Import

Config delivery settings can be imported using the id, e.g.

```terraform
terraform import tencentcloud_config_update_config_deliver.example example-id
```