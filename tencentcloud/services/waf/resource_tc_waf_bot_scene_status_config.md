Provides a resource to create a WAF bot scene status config

Example Usage

```hcl
resource "tencentcloud_waf_bot_scene_status_config" "example" {
  domain   = "example.com"
  scene_id = "3024324123"
  status   = true
}
```

Import

WAF bot scene status config can be imported using the id, e.g.

```
terraform import tencentcloud_waf_bot_scene_status_config.example example.com#3024324123
```
