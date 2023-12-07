Provides a resource to create a waf cc

Example Usage

```hcl
resource "tencentcloud_waf_cc" "example" {
  domain      = "www.demo.com"
  name        = "terraform"
  status      = 1
  advance     = "0"
  limit       = "60"
  interval    = "60"
  url         = "/cc_demo"
  match_func  = 0
  action_type = "22"
  priority    = 50
  valid_time  = 600
  edition     = "sparta-waf"
  type        = 1
}
```