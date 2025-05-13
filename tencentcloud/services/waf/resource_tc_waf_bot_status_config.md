Provides a resource to create a WAF bot status config

Example Usage

```hcl
resource "tencentcloud_waf_bot_status_config" "example" {
  domain   = "example.com"
  status   = "1"
}
```

Or

```hcl
resource "tencentcloud_waf_bot_status_config" "example" {
  domain      = "example.com"
  status      = "0"
  instance_id = "waf_2kxtlbky11bbcr4b"
}
```
