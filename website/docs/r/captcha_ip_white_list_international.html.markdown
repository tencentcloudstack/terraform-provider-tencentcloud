---
subcategory: "Captcha"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_captcha_ip_white_list_international"
sidebar_current: "docs-tencentcloud-resource-captcha_ip_white_list_international"
description: |-
  Provides a resource to create a captcha ip white list international.
---

# tencentcloud_captcha_ip_white_list_international

Provides a resource to create a captcha ip white list international.

## Example Usage

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

resource "tencentcloud_captcha_ip_white_list_international" "example" {
  name          = "tf-example"
  captcha_appid = tencentcloud_captcha_info_international.example.id
  ip            = "1.1.1.1"
  comment       = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `captcha_appid` - (Required, Int, ForceNew) Captcha appid.
* `ip` - (Required, String, ForceNew) IP data. Must be a single valid IPv4 or IPv6 address.
* `name` - (Required, String) IP allowlist name.
* `comment` - (Optional, String) Remark information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Captcha ip white list international can be imported using the id, e.g.

```
terraform import tencentcloud_captcha_ip_white_list_international.example 189910572#488584
```

