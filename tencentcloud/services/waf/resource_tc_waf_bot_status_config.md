Provides a resource to create a WAF bot status config

Example Usage

```hcl
resource "tencentcloud_waf_bot_status_config" "example" {
  instance_id = "waf_2kxtlbky11bbcr4b"
  domain      = "example.com"
  status      = "0"
}
```

Import

WAF bot status config can be imported using the id, e.g.

```
terraform import tencentcloud_waf_bot_status_config.example waf_2kxtlbky11bbcr4b#example.com
```
