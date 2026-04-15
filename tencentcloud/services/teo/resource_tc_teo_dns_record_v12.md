Provides a resource to create a TEO DNS record

Example Usage

### A Type DNS Record

```hcl
resource "tencentcloud_teo_dns_record_v12" "a_record" {
  zone_id = "zone-39quuimqg8r6"
  name    = "www"
  type    = "A"
  content = "1.2.3.4"
  ttl     = 300
}
```

### CNAME Type DNS Record with Weight

```hcl
resource "tencentcloud_teo_dns_record_v12" "cname_record" {
  zone_id = "zone-39quuimqg8r6"
  name    = "alias"
  type    = "CNAME"
  content = "target.example.com"
  weight  = 50
  location = "overseas"
}
```

### MX Type DNS Record with Priority

```hcl
resource "tencentcloud_teo_dns_record_v12" "mx_record" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "@"
  type     = "MX"
  content  = "mail.example.com"
  priority = 10
}
```

### Advanced Configuration with Weight and Location

```hcl
resource "tencentcloud_teo_dns_record_v12" "advanced_record" {
  zone_id  = "zone-39quuimqg8r6"
  name     = "www"
  type     = "A"
  content  = "1.2.3.4"
  location = "mainland"
  ttl      = 600
  weight   = 80
}
```

### Complete Resource Lifecycle (Create, Update, Delete)

```hcl
# Create a DNS record
resource "tencentcloud_teo_dns_record_v12" "example" {
  zone_id = "zone-39quuimqg8r6"
  name    = "www"
  type    = "A"
  content = "1.2.3.4"
  ttl     = 300
}

# Update the DNS record (change content and TTL)
resource "tencentcloud_teo_dns_record_v12" "example" {
  zone_id = "zone-39quuimqg8r6"
  name    = "www"
  type    = "A"
  content = "5.6.7.8"  # Updated content
  ttl     = 600         # Updated TTL
}

# Delete the DNS record (remove this resource from your Terraform configuration)
```

Import

TEO DNS record can be imported using the composite ID (zone_id#record_id), e.g.

```
terraform import tencentcloud_teo_dns_record_v12.example zone-39quuimqg8r6#record-67890
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) ID of the site.
* `name` - (Required) DNS record name.
* `type` - (Required) DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.
* `content` - (Required) DNS record content.
* `location` - (Optional) DNS record resolution route. Default is Default.
* `ttl` - (Optional, Computed) DNS record cache time. Value range: 60~86400, unit: seconds. Default: 300.
* `weight` - (Optional, Computed) DNS record weight. Value range: -1~100, -1 means no weight, 0 means no resolution. Default: -1.
* `priority` - (Optional, Computed) MX record priority. Value range: 0~50. Default: 0. Only valid when Type is MX.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `record_id` - DNS record ID.
* `status` - DNS record resolution status. Valid values: enable, disable.
* `created_on` - DNS record creation time.
