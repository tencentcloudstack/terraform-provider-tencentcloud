---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_play_auth_key_config"
sidebar_current: "docs-tencentcloud-resource-css_play_auth_key_config"
description: |-
  Provides a resource to create a css play_auth_key_config
---

# tencentcloud_css_play_auth_key_config

Provides a resource to create a css play_auth_key_config

## Example Usage

```hcl
resource "tencentcloud_css_play_auth_key_config" "play_auth_key_config" {
  domain_name   = "your_play_domain_name"
  enable        = 1
  auth_key      = "testauthkey"
  auth_delta    = 3600
  auth_back_key = "testbackkey"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Domain Name.
* `auth_back_key` - (Optional, String) Alternate key for authentication. No transfer means that the current value is not modified.
* `auth_delta` - (Optional, Int) Valid time, unit: second. No transfer means that the current value is not modified.
* `auth_key` - (Optional, String) Authentication key. No transfer means that the current value is not modified.
* `enable` - (Optional, Int) Enable or not, 0: Close, 1: Enable. No transfer means that the current value is not modified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css play_auth_key_config can be imported using the id, e.g.

```
terraform import tencentcloud_css_play_auth_key_config.play_auth_key_config play_auth_key_config_id
```

