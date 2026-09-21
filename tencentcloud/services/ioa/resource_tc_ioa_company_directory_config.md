Provides a resource to create a IOA company directory config.

Example Usage

```hcl
resource "tencentcloud_ioa_company_directory_config" "example" {
  type                = "WeCom"
  name                = "tf-example"
  config              = "encrypted-hex-config-data"
  sync_enable         = true
  sync_policy         = "daily"
  sync_policy_params  = jsonencode({ "hour" = 2 })
  create_auth_config  = true
  display_on_login_page = true
  description         = "tf example description"
  scene               = "API"

  name_i18n {
    lang  = "zh-CN"
    value = "示例目录"
  }

  name_i18n {
    lang  = "en-US"
    value = "Example Directory"
  }
}
```

Import

IOA company directory config can be imported using the id, e.g.

```
terraform import tencentcloud_ioa_company_directory_config.example 123
```
