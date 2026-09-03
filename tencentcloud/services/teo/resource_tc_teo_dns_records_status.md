Provides a resource to manage TEO (EdgeOne) DNS records status config.

Example Usage

Enable a DNS record

```hcl
resource "tencentcloud_teo_dns_records_status" "example" {
  zone_id    = "zone-3edjdliiw3he"
  records_id = "record-1234567890"
  status     = "enable"
}
```

Disable a DNS record

```hcl
resource "tencentcloud_teo_dns_records_status" "example" {
  zone_id    = "zone-3edjdliiw3he"
  records_id = "record-1234567890"
  status     = "disable"
}
```

Import

TEO (EdgeOne) DNS records status config can be imported using the composite id `zone_id#records_id`, e.g.

```
terraform import tencentcloud_teo_dns_records_status.example zone-3edjdliiw3he#record-1234567890
```
