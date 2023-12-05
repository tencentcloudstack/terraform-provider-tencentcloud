Provides a resource to create a waf cc_session

Example Usage

```hcl
resource "tencentcloud_waf_cc_session" "example" {
  domain           = "www.demo.com"
  source           = "get"
  category         = "match"
  key_or_start_mat = "key_a=123"
  end_mat          = "&"
  start_offset     = "-1"
  end_offset       = "-1"
  edition          = "sparta-waf"
  session_name     = "terraformDemo"
}
```

Import

waf cc_session can be imported using the id, e.g.

```
terraform import tencentcloud_waf_cc_session.example www.demo.com#sparta-waf#2000000253
```