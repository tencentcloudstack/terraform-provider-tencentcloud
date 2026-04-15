# tencentcloud_teo_dns_record_v11

Provides a resource to manage a DNS record of TEO.

## Example Usage

### Basic Usage

```hcl
resource "tencentcloud_teo_dns_record_v11" "example" {
    zone_id = tencentcloud_teo_zone.example.id
    name     = "www"
    type     = "A"
    content  = "1.2.3.4"
}
```

### With Optional Fields

```hcl
resource "tencentcloud_teo_dns_record_v11" "example" {
    zone_id  = tencentcloud_teo_zone.example.id
    name     = "www"
    type     = "A"
    content  = "1.2.3.4"
    ttl      = 300
    weight   = 50
    location = "Default"
}
```

### With MX Priority

```hcl
resource "tencentcloud_teo_dns_record_v11" "mx_example" {
    zone_id  = tencentcloud_teo_zone.example.id
    name     = "@"
    type     = "MX"
    content  = "mail.example.com"
    priority = 10
}
```

### Update Example

```hcl
# Initial creation
resource "tencentcloud_teo_dns_record_v11" "example" {
    zone_id = tencentcloud_teo_zone.example.id
    name     = "www"
    type     = "A"
    content  = "1.2.3.4"
    ttl      = 300
}

# Update ttl and content
resource "tencentcloud_teo_dns_record_v11" "example" {
    zone_id = tencentcloud_teo_zone.example.id
    name     = "www"
    type     = "A"
    content  = "5.6.7.8"
    ttl      = 600
}
```

### Different Record Types

#### A Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "a_record" {
    zone_id = tencentcloud_teo_zone.example.id
    name     = "www"
    type     = "A"
    content  = "1.2.3.4"
    weight   = 50
    location = "Default"
}
```

#### AAAA Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "aaaa_record" {
    zone_id = tencentcloud_teo_zone.example.id
    name     = "www"
    type     = "AAAA"
    content  = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
    weight   = 50
    location = "Default"
}
```

#### CNAME Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "cname_record" {
    zone_id  = tencentcloud_teo_zone.example.id
    name     = "www"
    type     = "CNAME"
    content  = "example.com"
    weight   = 50
    location = "Default"
}
```

#### MX Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "mx_record" {
    zone_id  = tencentcloud_teo_zone.example.id
    name     = "@"
    type     = "MX"
    content  = "mail.example.com"
    priority = 10
}
```

#### TXT Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "txt_record" {
    zone_id = tencentcloud_teo_zone.example.id
    name     = "_domainkey"
    type     = "TXT"
    content  = "v=spf1 include:example.com ~all"
    ttl      = 300
}
```

#### NS Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "ns_record" {
    zone_id = tencentcloud_teo_zone.example.id
    name     = "sub"
    type     = "NS"
    content  = "ns1.example.com"
    ttl      = 300
}
```

#### CAA Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "caa_record" {
    zone_id = tencentcloud_teo_zone.example.id
    name     = "@"
    type     = "CAA"
    content  = "0 issue letsencrypt.org"
    ttl      = 300
}
```

#### SRV Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "srv_record" {
    zone_id = tencentcloud_teo_zone.example.id
    name     = "_sip._tcp"
    type     = "SRV"
    content  = "10 60 5060 sipserver.example.com"
    ttl      = 300
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) ID of the site related with the DNS record.
* `name` - (Required, ForceNew) DNS record name. If it is Chinese, Korean, or Japanese domain name, it needs to be converted to punycode before input.
* `type` - (Required, ForceNew) DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.
* `content` - (Required) DNS record content. If it is Chinese, Korean, or Japanese domain name, it needs to be converted to punycode before input.
* `ttl` - (Optional) Cache time. Value range: 60~86400, default is 300, unit: seconds.
* `weight` - (Optional) DNS record weight. Value range: -1~100, default is -1 (not set). 0 means the record will not resolve. Only applicable for A, AAAA, and CNAME record types.
* `priority` - (Optional) MX record priority. Value range: 0~50, default is 0. Only applicable for MX record type.
* `location` - (Optional) DNS record resolution line. Default is Default, indicating the default resolution line, representing all regions. Only applicable for A, AAAA, and CNAME record types.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the DNS record. The format is `zoneId#recordId`.
* `status` - DNS record resolution status. Valid values: enable, disable.
* `created_on` - Creation time of the DNS record.
* `modified_on` - Modification time of the DNS record.

## Import

DNS record can be imported, e.g.

```bash
$ terraform import tencentcloud_teo_dns_record_v11.example zoneId#recordId
```
