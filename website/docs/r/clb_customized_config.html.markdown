---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_customized_config"
sidebar_current: "docs-tencentcloud-resource-clb_customized_config"
description: |-
  Provides a resource to create a CLB customized config which type is `CLB`.
---

# tencentcloud_clb_customized_config

Provides a resource to create a CLB customized config which type is `CLB`.

## Example Usage

### Create clb customized config without CLB instance

```hcl
resource "tencentcloud_clb_customized_config" "example" {
  config_name    = "tf-example"
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
}
```

### Create clb customized config with CLB instances

```hcl
resource "tencentcloud_clb_customized_config" "example" {
  config_name    = "tf-example"
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  load_balancer_ids = [
    "lb-l6cp6jt4",
    "lb-muk4zzxi",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `config_content` - (Required, String) Content of Customized Config.
* `config_name` - (Required, String) Name of Customized Config.
* `load_balancer_ids` - (Optional, Set: [`String`]) List of LoadBalancer Ids.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of Customized Config.
* `update_time` - Update time of Customized Config.


## Import

CLB customized config can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_customized_config.example pz-diowqstq
```

