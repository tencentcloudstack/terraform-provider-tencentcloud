Provide a resource to create a Private Dns Record.

Example Usage

```hcl
resource "tencentcloud_private_dns_record" "foo" {
  zone_id      = "zone-rqndjnki"
  record_type  = "A"
  record_value = "192.168.1.2"
  sub_domain   = "www"
  ttl          = 300
  weight       = 1
  mx           = 0
}
```

Import

Private Dns Record can be imported, e.g.

```
$ terraform import tencentcloud_private_dns_zone.foo zone_id#record_id
```