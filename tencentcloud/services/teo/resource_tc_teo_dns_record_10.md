# tencentcloud_teo_dns_record_10

Provides a TEO DNS Record resource.

## Example Usage

### Basic A Record

```hcl
resource "tencentcloud_teo_dns_record_10" "example_a" {
  zone_id = "zone-12345678"
  name     = "www"
  type     = "A"
  content  = "1.2.3.4"
}
```

### MX Record with Priority

```hcl
resource "tencentcloud_teo_dns_record_10" "example_mx" {
  zone_id  = "zone-12345678"
  name      = "@"
  type      = "MX"
  content   = "mail.example.com"
  priority  = 10
}
```

### CNAME Record

```hcl
resource "tencentcloud_teo_dns_record_10" "example_cname" {
  zone_id = "zone-12345678"
  name     = "www"
  type     = "CNAME"
  content  = "another-domain.com"
}
```

### Record with Optional Parameters

```hcl
resource "tencentcloud_teo_dns_record_10" "example_optional" {
  zone_id  = "zone-12345678"
  name      = "www"
  type      = "A"
  content   = "1.2.3.4"
  location  = "China"
  ttl       = 600
  weight    = 50
}
```

### TXT Record for SPF

```hcl
resource "tencentcloud_teo_dns_record_10" "example_txt" {
  zone_id = "zone-12345678"
  name     = "@"
  type     = "TXT"
  content  = "v=spf1 include:_spf.example.com ~all"
}
```

### AAAA Record

```hcl
resource "tencentcloud_teo_dns_record_10" "example_aaaa" {
  zone_id = "zone-12345678"
  name     = "www"
  type     = "AAAA"
  content  = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) ID of the site related with the DNS record.
* `name` - (Required) DNS record name. If it is a Chinese, Korean, or Japanese domain name, it needs to be converted to punycode before input.
* `type` - (Required) DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.
* `content` - (Required) DNS record content. Fill in the corresponding content according to the Type value. If it is a Chinese, Korean, or Japanese domain name, it needs to be converted to punycode before input.
* `location` - (Optional, Computed) DNS record resolution line. Default is Default, which means default resolution line and represents all regions. Valid only when Type is A, AAAA, or CNAME.
* `ttl` - (Optional, Computed) Cache time. Valid range: 60-86400, unit: seconds. Default is 300.
* `weight` - (Optional, Computed) DNS record weight. Valid range: -1-100. Default is -1, which means no weight is set. When set to 0, it means no resolution. Valid only when Type is A, AAAA, or CNAME.
* `priority` - (Optional, Computed) MX record priority. Valid only when Type is MX. Valid range: 0-50, smaller value, higher priority. Default is 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `record_id` - DNS record ID.
* `status` - DNS record resolution status. Valid values: enable, disable.
* `created_on` - Creation time of the DNS record.

## DNS Record Types

### A Record
Maps a domain name to an IPv4 address.

Example:
```hcl
resource "tencentcloud_teo_dns_record_10" "a_record" {
  zone_id = "zone-12345678"
  name     = "www"
  type     = "A"
  content  = "192.0.2.1"
}
```

### AAAA Record
Maps a domain name to an IPv6 address.

Example:
```hcl
resource "tencentcloud_teo_dns_record_10" "aaaa_record" {
  zone_id = "zone-12345678"
  name     = "www"
  type     = "AAAA"
  content  = "2001:db8::1"
}
```

### MX Record
Specifies the mail server for a domain. The priority value determines which server is preferred.

Example:
```hcl
resource "tencentcloud_teo_dns_record_10" "mx_record" {
  zone_id  = "zone-12345678"
  name      = "@"
  type      = "MX"
  content   = "mail.example.com"
  priority  = 10
}
```

### CNAME Record
Maps a domain name to another domain name.

Example:
```hcl
resource "tencentcloud_teo_dns_record_10" "cname_record" {
  zone_id = "zone-12345678"
  name     = "www"
  type     = "CNAME"
  content  = "another-domain.com"
}
```

### TXT Record
Stores text information about the domain. Commonly used for SPF records, DKIM, domain verification, etc.

Example:
```hcl
resource "tencentcloud_teo_dns_record_10" "txt_record" {
  zone_id = "zone-12345678"
  name     = "_acme-challenge"
  type     = "TXT"
  content  = "verification_token_here"
}
```

### NS Record
Specifies the name server for the domain.

Example:
```hcl
resource "tencentcloud_teo_dns_record_10" "ns_record" {
  zone_id = "zone-12345678"
  name     = "subdomain"
  type     = "NS"
  content  = "ns1.example.com"
}
```

### CAA Record
Specifies which certificate authorities (CAs) are allowed to issue certificates for the domain.

Example:
```hcl
resource "tencentcloud_teo_dns_record_10" "caa_record" {
  zone_id = "zone-12345678"
  name     = "@"
  type     = "CAA"
  content  = "0 issue letsencrypt.org"
}
```

### SRV Record
Specifies the location of a specific service.

Example:
```hcl
resource "tencentcloud_teo_dns_record_10" "srv_record" {
  zone_id = "zone-12345678"
  name     = "_sip._tcp"
  type     = "SRV"
  content  = "10 60 5060 sipserver.example.com"
}
```

## Parameter Constraints

### TTL
- Valid range: 60-86400 seconds
- Default: 300 seconds
- Smaller value means faster propagation but more queries

### Weight
- Valid range: -1-100
- Default: -1 (no weight set)
- Value 0 means the record will not resolve
- Only applicable to A, AAAA, and CNAME record types
- Used for load balancing across multiple records with the same name

### Priority
- Valid range: 0-50
- Default: 0
- Smaller value = higher priority
- Only applicable to MX record type
- Used when multiple MX records exist for the same domain

### Location
- Default: "Default" (all regions)
- Only applicable to A, AAAA, and CNAME record types
- Used for geo-based DNS routing

## Import

DNS record can be imported using the composite ID in the format `zoneId#recordId`:

```shell
terraform import tencentcloud_teo_dns_record_10.example zone-12345678#record-87654321
```

## Timeouts

The `timeouts` block allows you to specify timeouts for certain operations:

* `create` - (Default 10 minutes) Time to wait for DNS record creation
* `read` - (Default 3 minutes) Time to wait for DNS record read
* `update` - (Default 10 minutes) Time to wait for DNS record update
* `delete` - (Default 10 minutes) Time to wait for DNS record deletion

Example:
```hcl
resource "tencentcloud_teo_dns_record_10" "example" {
  zone_id = "zone-12345678"
  name     = "www"
  type     = "A"
  content  = "1.2.3.4"

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
```

## Notes

- DNS record operations are asynchronous. The resource will wait for the record to be propagated before completing.
- The resource supports idempotent creation. If a record with the same ZoneId, Name, Type, and Content already exists, the existing record will be used instead of creating a duplicate.
- When modifying DNS records, only the changed fields will be updated. No changes detected will result in no API call being made.
- The resource ID is in the format `zoneId#recordId` for easy identification and import.
