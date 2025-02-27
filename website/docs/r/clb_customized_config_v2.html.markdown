---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_customized_config_v2"
sidebar_current: "docs-tencentcloud-resource-clb_customized_config_v2"
description: |-
  Provides a resource to create a CLB customized config which type is `SERVER` or `LOCATION`.
---

# tencentcloud_clb_customized_config_v2

Provides a resource to create a CLB customized config which type is `SERVER` or `LOCATION`.

## Example Usage

### If config_type is SERVER

```hcl
resource "tencentcloud_clb_customized_config_v2" "example" {
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  config_name    = "tf-example"
  config_type    = "SERVER"
}

output "configId" {
  value = tencentcloud_clb_customized_config_v2.example.config_id
}
```

### If config_type is LOCATION

```hcl
resource "tencentcloud_clb_customized_config_v2" "example" {
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  config_name    = "tf-example"
  config_type    = "LOCATION"
}

output "configId" {
  value = tencentcloud_clb_customized_config_v2.example.config_id
}
```

## Argument Reference

The following arguments are supported:

* `config_content` - (Required, String) Content of Customized Config.
* `config_name` - (Required, String) Name of Customized Config.
* `config_type` - (Required, String, ForceNew) Type of Customized Config. Valid values: `SERVER` and `LOCATION`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `config_id` - ID of Customized Config.
* `create_time` - Create time of Customized Config.
* `update_time` - Update time of Customized Config.


## Import

CLB customized V2 config can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_customized_config_v2.example pz-diowqstq#SERVER

Or

$ terraform import tencentcloud_clb_customized_config_v2.example pz-4r10y4b2#LOCATION
```

