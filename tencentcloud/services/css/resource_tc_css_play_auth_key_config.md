Provides a resource to create a css play_auth_key_config

Example Usage

```hcl
resource "tencentcloud_css_play_auth_key_config" "play_auth_key_config" {
  domain_name = "your_play_domain_name"
  enable = 1
  auth_key = "testauthkey"
  auth_delta = 3600
  auth_back_key = "testbackkey"
}
```

Import

css play_auth_key_config can be imported using the id, e.g.

```
terraform import tencentcloud_css_play_auth_key_config.play_auth_key_config play_auth_key_config_id
```