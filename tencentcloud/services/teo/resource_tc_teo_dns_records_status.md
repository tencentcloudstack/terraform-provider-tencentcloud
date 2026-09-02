Provides a resource to manage TEO (EdgeOne) DNS records status config.

Example Usage

Enable a DNS record

```hcl
resource "tencentcloud_teo_dns_records_status" "example" {
  zone_id = "zone-3edjdliiw3he"

  filters {
    name   = "id"
    values = ["record-1234567890"]
  }

  records_to_enable = ["record-1234567890"]
}
```

Disable a DNS record

```hcl
resource "tencentcloud_teo_dns_records_status" "example" {
  zone_id = "zone-3edjdliiw3he"

  filters {
    name   = "id"
    values = ["record-1234567890"]
  }

  records_to_disable = ["record-1234567890"]
}
```

Import

TEO (EdgeOne) DNS records status config can be imported using the `zone_id`, e.g.

```
terraform import tencentcloud_teo_dns_records_status.example zone-3edjdliiw3he
```
