Provides a resource to create a waf ip_access_control

Example Usage

```hcl
resource "tencentcloud_waf_ip_access_control" "example" {
  instance_id = "waf_2kxtlbky00b3b4qz"
  domain      = "www.demo.com"
  edition     = "sparta-waf"
  items {
    ip       = "1.1.1.1"
    note     = "desc info."
    action   = 40
    valid_ts = "2019571199"
  }

  items {
    ip       = "2.2.2.2"
    note     = "desc info."
    action   = 42
    valid_ts = "2019571199"
  }

  items {
    ip       = "3.3.3.3"
    note     = "desc info."
    action   = 40
    valid_ts = "1680570420"
  }
}
```

Import

waf ip_access_control can be imported using the id, e.g.

```
terraform import tencentcloud_waf_ip_access_control.example waf_2kxtlbky00b3b4qz#www.demo.com#sparta-waf
```