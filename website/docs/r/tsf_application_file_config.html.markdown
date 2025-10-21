---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_application_file_config"
sidebar_current: "docs-tencentcloud-resource-tsf_application_file_config"
description: |-
  Provides a resource to create a tsf application_file_config
---

# tencentcloud_tsf_application_file_config

Provides a resource to create a tsf application_file_config

## Example Usage

```hcl
resource "tencentcloud_tsf_application_file_config" "application_file_config" {
  config_name         = "terraform-test"
  config_version      = "1.0"
  config_file_name    = "application.yaml"
  config_file_value   = "test: 1"
  application_id      = "application-a24x29xv"
  config_file_path    = "/etc/nginx"
  config_version_desc = "1.0"
  config_file_code    = "UTF-8"
  config_post_cmd     = "source .bashrc"
  encode_with_base64  = true
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required, String, ForceNew) Config file associated application ID.
* `config_file_name` - (Required, String, ForceNew) Config file name.
* `config_file_path` - (Required, String, ForceNew) config release path.
* `config_file_value` - (Required, String, ForceNew) Configuration file content (the original content encoding needs to be in utf-8 format, if the ConfigFileCode is gbk, it will be converted in the background).
* `config_name` - (Required, String, ForceNew) Config Name.
* `config_version` - (Required, String, ForceNew) Config version.
* `config_file_code` - (Optional, String, ForceNew) Configuration file encoding, utf-8 or gbk. Note: If you choose gbk, you need the support of a new version of tsf-consul-template (public cloud virtual machines need to use 1.32 tsf-agent, and containers need to obtain the latest tsf-consul-template-docker.tar.gz from the documentation).
* `config_post_cmd` - (Optional, String, ForceNew) post command.
* `config_version_desc` - (Optional, String, ForceNew) config version description.
* `encode_with_base64` - (Optional, Bool, ForceNew) the config value is encoded with base64 or not.
* `program_id_list` - (Optional, Set: [`String`], ForceNew) datasource for auth.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



