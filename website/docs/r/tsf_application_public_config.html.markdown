---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_application_public_config"
sidebar_current: "docs-tencentcloud-resource-tsf_application_public_config"
description: |-
  Provides a resource to create a tsf application_public_config
---

# tencentcloud_tsf_application_public_config

Provides a resource to create a tsf application_public_config

## Example Usage

```hcl
resource "tencentcloud_tsf_application_public_config" "application_public_config" {
  config_name         = "my_config"
  config_version      = "1.0"
  config_value        = "test: 1"
  config_version_desc = "product version"
  config_type         = "P"
  encode_with_base64  = true
  # program_id_list =
}
```

## Argument Reference

The following arguments are supported:

* `config_name` - (Required, String, ForceNew) Config Name.
* `config_value` - (Required, String, ForceNew) config value, only yaml file allowed.
* `config_version` - (Required, String, ForceNew) config version.
* `config_type` - (Optional, String, ForceNew) Config type.
* `config_version_desc` - (Optional, String, ForceNew) Config version description.
* `encode_with_base64` - (Optional, Bool, ForceNew) the config value is encoded with base64 or not.
* `program_id_list` - (Optional, Set: [`String`], ForceNew) datasource for auth.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



