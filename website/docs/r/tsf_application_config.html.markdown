---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_application_config"
sidebar_current: "docs-tencentcloud-resource-tsf_application_config"
description: |-
  Provides a resource to create a tsf application_config
---

# tencentcloud_tsf_application_config

Provides a resource to create a tsf application_config

## Example Usage

```hcl
resource "tencentcloud_tsf_application_config" "application_config" {
  config_name         = "test-2"
  config_version      = "1.0"
  config_value        = "name: \"name\""
  application_id      = "application-ym9mxmza"
  config_version_desc = "test2"
  # config_type = ""
  encode_with_base64 = false
  # program_id_list =
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required, String) Application ID.
* `config_name` - (Required, String) configuration item name.
* `config_value` - (Required, String) configuration item value.
* `config_version` - (Required, String) configuration item version.
* `config_type` - (Optional, String) configuration item value type.
* `config_version_desc` - (Optional, String) configuration item version description.
* `encode_with_base64` - (Optional, Bool) Base64 encoded configuration items.
* `program_id_list` - (Optional, Set: [`String`]) Program id list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tsf application_config can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_config.application_config dcfg-y4e3zngv
```

