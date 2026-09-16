---
subcategory: "Identity and Access Management Open APIs(IOA)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ioa_company_directory_config"
sidebar_current: "docs-tencentcloud-resource-ioa_company_directory_config"
description: |-
  Provides a resource to create a IOA company directory config.
---

# tencentcloud_ioa_company_directory_config

Provides a resource to create a IOA company directory config.

## Example Usage

```hcl
resource "tencentcloud_ioa_company_directory_config" "example" {
  type                  = "WeCom"
  name                  = "tf-example"
  config                = "encrypted-hex-config-data"
  sync_enable           = true
  sync_policy           = "daily"
  sync_policy_params    = jsonencode({ "hour" = 2 })
  create_auth_config    = true
  display_on_login_page = true
  description           = "tf example description"
  scene                 = "API"

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

## Argument Reference

The following arguments are supported:

* `config` - (Required, String) Configuration data encrypted by SM2 then Hex encoded.
* `create_auth_config` - (Required, Bool) Whether to synchronously create an auth source.
* `display_on_login_page` - (Required, Bool) Whether to display on the login page.
* `name` - (Required, String) Enterprise directory name.
* `sync_enable` - (Required, Bool) Whether scheduled sync is enabled.
* `sync_policy_params` - (Required, String) JSON string of sync strategy parameters.
* `sync_policy` - (Required, String) Scheduled sync strategy, enum: `4hours`, `daily`, `weekly`.
* `type` - (Required, String) Enterprise directory type, enum: `WeCom`, `Lark`, `DingTalk`, `MicrosoftEntraID`.
* `description` - (Optional, String) Directory description.
* `name_i18n` - (Optional, List) Multi-language names.
* `scene` - (Optional, String) Usage scene: API created, quick start, normal config. Create-only, not present in Describe output.

The `name_i18n` object supports the following:

* `lang` - (Required, String) Language enum, e.g. `zh-CN`, `en-US`.
* `value` - (Required, String) Localized name value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

IOA company directory config can be imported using the id, e.g.

```
terraform import tencentcloud_ioa_company_directory_config.example 123
```

