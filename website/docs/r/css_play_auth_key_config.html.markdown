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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyLivePlayAuthKey](https://www.tencentcloud.com/document/product/267/30770) - This API is used to modify the playback authentication key.
* Read: [DescribeLivePlayAuthKey](https://www.tencentcloud.com/document/product/267/30772) - This API is used to query the playback authentication key.
* Update: [ModifyLivePlayAuthKey](https://www.tencentcloud.com/document/product/267/30770) - This API is used to modify the playback authentication key.

## Import

css play_auth_key_config can be imported using the id, e.g.

```
terraform import tencentcloud_css_play_auth_key_config.play_auth_key_config play_auth_key_config_id
```

