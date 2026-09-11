---
subcategory: "Captcha"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_captcha_info_international"
sidebar_current: "docs-tencentcloud-resource-captcha_info_international"
description: |-
  Provides a resource to create a captcha info international.
---

# tencentcloud_captcha_info_international

Provides a resource to create a captcha info international.

## Example Usage

### If channel_info is web

```hcl
resource "tencentcloud_captcha_info_international" "example" {
  app_name                 = "tf-example"
  channel_info             = "web"
  verify_rank              = "1"
  user_set_cap_type        = "2"
  defend_mode              = "block"
  disable_invisible_switch = "0"
  verify_domain            = "example.com"
  check_appid_switch       = 1
  check_iv_switch          = 1
  check_box_style          = "2"
}
```

### If channel_info is ios

```hcl
resource "tencentcloud_captcha_info_international" "example" {
  app_name                 = "tf-example"
  channel_info             = "ios"
  verify_rank              = "3"
  user_set_cap_type        = "1"
  defend_mode              = "notify"
  disable_invisible_switch = "2"
  verify_bundle_id         = "com.example.SampleApp"
  check_appid_switch       = 0
  check_iv_switch          = 0
  check_box_style          = "2"
}
```

### If channel_info is android

```hcl
resource "tencentcloud_captcha_info_international" "example" {
  app_name                 = "tf-example"
  channel_info             = "android"
  verify_rank              = "2"
  user_set_cap_type        = "2"
  defend_mode              = "block"
  disable_invisible_switch = "1"
  verify_package           = "com.example.myapp"
  check_appid_switch       = 0
  check_iv_switch          = 0
  check_box_style          = "2"
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required, String) Captcha name.
* `channel_info` - (Required, String, ForceNew) Client type. Valid values: `web` (web scenario), `android` (android client), `ios` (ios client). Default is `web`.
* `check_appid_switch` - (Optional, Int) Whether to enable captcha encryption. `0`: disable, `1`: enable.
* `check_box_style` - (Optional, String) Checkbox style. Valid values: `0` (simple), `1` (basic), `2` (invisible).
* `check_iv_switch` - (Optional, Int) Whether to enable one-time pad. Valid values: `0` (disable), `1` (enable). Only allowed to be `1` when `check_appid_switch` is `1`.
* `defend_mode` - (Optional, String) Defend mode. Valid values: `block` (block mode), `notify` (perceive mode). Default is `notify`.
* `disable_invisible_switch` - (Optional, String) Verify mechanism. Valid values: `0` (one-click verify), `1` (always verify), `2` (invisible verify, UserSetCapType must be `1`).
* `tags` - (Optional, List: [`String`]) Resource tags, key&value format.
* `user_set_cap_type` - (Optional, String) Verify type. Valid values: `1` (invisible verify, DisableInvisibleSwitch must be `2`), `2` (slider), `8` (image click), `9` (voice).
* `verify_bundle_id` - (Optional, String) App BundleId. Only valid when `channel_info` is `ios`.
* `verify_domain` - (Optional, String) Web domain. Only valid when `channel_info` is `web`.
* `verify_package` - (Optional, String) App package. Only valid when `channel_info` is `android`.
* `verify_rank` - (Optional, String) Verify rank. Valid values: `1` (experience first), `2` (balanced), `3` (security first). Default is `1`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Captcha info international can be imported using the id, e.g.

```
terraform import tencentcloud_captcha_info_international.example 189933256
```

