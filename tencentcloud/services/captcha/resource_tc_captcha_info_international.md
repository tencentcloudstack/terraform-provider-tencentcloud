Provides a resource to create a captcha info international.

Example Usage

If channel_info is web

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

If channel_info is ios

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

If channel_info is android

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

Import

Captcha info international can be imported using the id, e.g.

```
terraform import tencentcloud_captcha_info_international.example 189933256
```
