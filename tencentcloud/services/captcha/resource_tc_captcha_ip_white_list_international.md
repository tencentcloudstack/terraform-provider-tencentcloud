Provides a resource to create a captcha ip white list international.

Example Usage

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

Import

Captcha ip white list international can be imported using the id, e.g.

```
terraform import tencentcloud_captcha_ip_white_list_international.example 189910572#488584
```
