Provides a resource to set clb listener default domain

Example Usage

Set default domain

```hcl
resource "tencentcloud_clb_listener_default_domain" "example" {
  clb_id      = "lb-g1miv1ok"
  listener_id = "lbl-duilx5qm"
  domain      = "3.com"
}
```

Import

CLB listener default domain can be imported using the id (version >= 1.47.0), e.g.

```
$ terraform import tencentcloud_clb_listener_default_domain.example lb-k2zjp9lv#lbl-hh141sn9
```
