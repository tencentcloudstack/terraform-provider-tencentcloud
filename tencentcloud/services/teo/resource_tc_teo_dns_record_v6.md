Provides a resource to create a teo teo_dns_record_v6

Example Usage

### Basic Usage - A Record

```hcl
resource "tencentcloud_teo_dns_record_v6" "a_record" {
  zone_id  = "zone-xxxxxxxx"
  type     = "A"
  content  = "1.2.3.4"
  name     = "www.example.com"
  ttl      = 300
}
```

### AAAA Record (IPv6)

```hcl
resource "tencentcloud_teo_dns_record_v6" "aaaa_record" {
  zone_id  = "zone-xxxxxxxx"
  type     = "AAAA"
  content  = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
  name     = "www.example.com"
  ttl      = 600
}
```

### CNAME Record

```hcl
resource "tencentcloud_teo_dns_record_v6" "cname_record" {
  zone_id  = "zone-xxxxxxxx"
  type     = "CNAME"
  content  = "example.com"
  name     = "www.example.com"
  location = "Default"
  weight   = 10
  ttl      = 300
}
```

### MX Record with Priority

```hcl
resource "tencentcloud_teo_dns_record_v6" "mx_record" {
  zone_id  = "zone-xxxxxxxx"
  type     = "MX"
  content  = "mail.example.com"
  name     = "example.com"
  priority = 10
  ttl      = 300
}
```

### TXT Record

```hcl
resource "tencentcloud_teo_dns_record_v6" "txt_record" {
  zone_id  = "zone-xxxxxxxx"
  type     = "TXT"
  content  = "v=spf1 include:_spf.example.com ~all"
  name     = "example.com"
  ttl      = 300
}
```

### NS Record

```hcl
resource "tencentcloud_teo_dns_record_v6" "ns_record" {
  zone_id  = "zone-xxxxxxxx"
  type     = "NS"
  content  = "ns1.example.com"
  name     = "subdomain.example.com"
  ttl      = 300
}
```

### CAA Record

```hcl
resource "tencentcloud_teo_dns_record_v6" "caa_record" {
  zone_id  = "zone-xxxxxxxx"
  type     = "CAA"
  content  = "0 issue \"letsencrypt.org\""
  name     = "example.com"
  ttl      = 300
}
```

### SRV Record

```hcl
resource "tencentcloud_teo_dns_record_v6" "srv_record" {
  zone_id  = "zone-xxxxxxxx"
  type     = "SRV"
  content  = "10 60 5060 sipserver.example.com"
  name     = "_sip._tcp.example.com"
  ttl      = 300
}
```

### Complete Example with All Parameters

```hcl
resource "tencentcloud_teo_dns_record_v6" "complete" {
  zone_id  = "zone-xxxxxxxx"
  type     = "A"
  content  = "1.2.3.4"
  name     = "www.example.com"
  location = "Default"
  ttl      = 300
  weight   = 50
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Zone id.
* `name` - (Required) DNS record name. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.
* `type` - (Required) DNS record type. valid values are:
  - A: points the domain name to an external ipv4 address, such as 8.8.8.8
  - AAAA: points the domain name to an external ipv6 address
  - MX: used for email servers. when there are multiple mx records, the lower the priority value, the higher the priority
  - CNAME: points the domain name to another domain name, which then resolves to the final ip address
  - TXT: identifies and describes the domain name, commonly used for domain verification and spf records (anti-spam)
  - NS: if you need to delegate the subdomain to another dns service provider for resolution, you need to add an ns record. the root domain cannot add ns records
  - CAA: specifies the ca that can issue certificates for this site
  - SRV: identifies a server using a service, commonly used in microsoft's directory management

  Different record types, such as SRV and CAA records, have different requirements for host record names and record value formats. for detailed descriptions and format examples of each record type, please refer to: [introduction to dns record types](https://intl.cloud.tencent.com/document/product/1552/90453?from_cn_redirect=1#2f681022-91ab-4a9e-ac3d-0a6c454d954e).
* `content` - (Required) DNS record content. fill in the corresponding content according to the type value. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.
* `location` - (Optional) DNS record resolution route. if not specified, the default is DEFAULT, which means the default resolution route and is effective in all regions.
  - resolution route configuration is only applicable when type (dns record type) is A, AAAA, or CNAME.
  - resolution route configuration is only applicable to standard version and enterprise edition packages. for valid values, please refer to: [resolution routes and corresponding code enumeration](https://intl.cloud.tencent.com/document/product/1552/112542?from_cn_redirect=1).
* `ttl` - (Optional, Default: `300`) Cache time. users can specify a value range of 60-86400. the smaller the value, the faster the modification records will take effect in all regions. unit: seconds.
* `weight` - (Optional, Default: `-1`) DNS record weight. users can specify a value range of -1 to 100. a value of 0 means no resolution. if not specified, the default is -1, which means no weight is set. weight configuration is only applicable when type (dns record type) is A, AAAA, or CNAME. note: for the same subdomain, different dns records with the same resolution route should either all have weights set or none have weights set.
* `priority` - (Optional, Default: `0`) MX record priority, which takes effect only when type (dns record type) is MX. the smaller the value, the higher the priority. users can specify a value range of 0-50.
* `status` - (Optional, Computed) DNS record resolution status, the following values:
  - enable: has taken effect
  - disable: has been disabled

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in the format of `zoneId#recordId`.
* `record_id` - DNS record ID.
* `created_on` - Creation time.
* `modified_on` - Modify time.

## Parameter Limitations

### TTL
- **Valid Range**: 60-86400 seconds
- **Default**: 300 seconds
- **Description**: The smaller the value, the faster the DNS record changes take effect globally

### Weight
- **Valid Range**: -1 to 100
- **Default**: -1 (no weight set)
- **Special Values**:
  - -1: No weight set (default)
  - 0: No resolution (record is disabled for this route)
  - 1-100: Weight value used for weighted routing
- **Applicable Types**: A, AAAA, CNAME only
- **Note**: For the same subdomain, different DNS records with the same resolution route should either all have weights set or none have weights set

### Priority
- **Valid Range**: 0 to 50
- **Default**: 0
- **Description**: The smaller the value, the higher the priority
- **Applicable Types**: MX only

### Location
- **Default**: Default (effective in all regions)
- **Applicable Types**: A, AAAA, CNAME only
- **Note**: Only applicable to standard version and enterprise edition packages

## DNS Record Type Requirements

### A Record
- Points the domain name to an external IPv4 address
- Supports location and weight parameters
- Example content: `1.2.3.4`

### AAAA Record
- Points the domain name to an external IPv6 address
- Supports location and weight parameters
- Example content: `2001:0db8:85a3:0000:0000:8a2e:0370:7334`

### CNAME Record
- Points the domain name to another domain name
- Supports location and weight parameters
- Example content: `example.com`
- Note: CNAME records cannot coexist with other record types for the same host

### MX Record
- Used for email servers
- Supports priority parameter
- Example content: `mail.example.com`
- Note: Lower priority value = higher priority

### TXT Record
- Used for domain verification and SPF records
- Example content: `v=spf1 include:_spf.example.com ~all`
- Note: Text content should be enclosed in quotes

### NS Record
- Used to delegate subdomain to another DNS service provider
- Example content: `ns1.example.com`
- Note: Root domain cannot add NS records

### CAA Record
- Specifies the CA that can issue certificates for this site
- Example content: `0 issue "letsencrypt.org"`
- Format: `[flags] [tag] [value]`

### SRV Record
- Identifies a server using a service
- Example content: `10 60 5060 sipserver.example.com`
- Format: `[priority] [weight] [port] [target]`

## International Domain Names (IDN)

Chinese, Korean, or Japanese domain names need to be converted to punycode format before input.

Example:
- Original: `中文.com`
- Punycode: `xn--fiq228c.com`

Use an IDN to punycode converter if needed.

## Import

teo teo_dns_record_v6 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_v6.example {zoneId}#{recordId}
```

Example:
```
terraform import tencentcloud_teo_dns_record_v6.example zone-123abc456def#record-789xyz012uvw
```
