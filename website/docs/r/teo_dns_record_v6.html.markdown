---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_v6"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_v6"
description: |-
  Provides a resource to create a teo teo_dns_record_v6
---

# tencentcloud_teo_dns_record_v6

Provides a resource to create a teo teo_dns_record_v6

## Example Usage

### ### Basic Usage - A Record

```hcl
resource "tencentcloud_teo_dns_record_v6" "a_record" {
  zone_id = "zone-xxxxxxxx"
  type    = "A"
  content = "1.2.3.4"
  name    = "www.example.com"
  ttl     = 300
}
```

### ### AAAA Record (IPv6)

```hcl
resource "tencentcloud_teo_dns_record_v6" "aaaa_record" {
  zone_id = "zone-xxxxxxxx"
  type    = "AAAA"
  content = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
  name    = "www.example.com"
  ttl     = 600
}
```

### ### CNAME Record

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

### ### MX Record with Priority

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

### ### TXT Record

```hcl
resource "tencentcloud_teo_dns_record_v6" "txt_record" {
  zone_id = "zone-xxxxxxxx"
  type    = "TXT"
  content = "v=spf1 include:_spf.example.com ~all"
  name    = "example.com"
  ttl     = 300
}
```

### ### NS Record

```hcl
resource "tencentcloud_teo_dns_record_v6" "ns_record" {
  zone_id = "zone-xxxxxxxx"
  type    = "NS"
  content = "ns1.example.com"
  name    = "subdomain.example.com"
  ttl     = 300
}
```

### ### CAA Record

```hcl
resource "tencentcloud_teo_dns_record_v6" "caa_record" {
  zone_id = "zone-xxxxxxxx"
  type    = "CAA"
  content = "0 issue \"letsencrypt.org\""
  name    = "example.com"
  ttl     = 300
}
```

### ### SRV Record

```hcl
resource "tencentcloud_teo_dns_record_v6" "srv_record" {
  zone_id = "zone-xxxxxxxx"
  type    = "SRV"
  content = "10 60 5060 sipserver.example.com"
  name    = "_sip._tcp.example.com"
  ttl     = 300
}
```

### ### Complete Example with All Parameters

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

### Use an IDN to punycode converter if needed.

## Import

teo teo_dns_record_v6 can be imported using the id, e.g.

```hcl
terraform import tencentcloud_teo_dns_record_v6.example { zoneId } #{recordId}
```

### Example:

```hcl
terraform import tencentcloud_teo_dns_record_v6.example zone-123abc456def #record-789xyz012uvw
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record content. fill in the corresponding content according to the type value. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.
* `name` - (Required, String) DNS record name. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.
* `type` - (Required, String) DNS record type. valid values are:
	- A: points the domain name to an external ipv4 address, such as 8.8.8.8;
	- AAAA: points the domain name to an external ipv6 address;
	- MX: used for email servers. when there are multiple mx records, the lower the priority value, the higher the priority;
	- CNAME: points the domain name to another domain name, which then resolves to the final ip address;
	- TXT: identifies and describes the domain name, commonly used for domain verification and spf records (anti-spam);
	- NS: if you need to delegate the subdomain to another dns service provider for resolution, you need to add an ns record. the root domain cannot add ns records;
	- CAA: specifies the ca that can issue certificates for this site;
	- SRV: identifies a server using a service, commonly used in microsoft's directory management.
Different record types, such as SRV and CAA records, have different requirements for host record names and record value formats. for detailed descriptions and format examples of each record type, please refer to: [introduction to dns record types](https://intl.cloud.tencent.com/document/product/1552/90453?from_cn_redirect=1#2f681022-91ab-4a9e-ac3d-0a6c454d954e).
* `zone_id` - (Required, String, ForceNew) Zone id.
* `location` - (Optional, String) DNS record resolution route. if not specified, the default is DEFAULT, which means the default resolution route and is effective in all regions.

- resolution route configuration is only applicable when type (dns record type) is A, AAAA, or CNAME.
- resolution route configuration is only applicable to standard version and enterprise edition packages. for valid values, please refer to: [resolution routes and corresponding code enumeration](https://intl.cloud.tencent.com/document/product/1552/112542?from_cn_redirect=1).
* `priority` - (Optional, Int) MX record priority, which takes effect only when type (dns record type) is MX. the smaller the value, the higher the priority. users can specify a value range of 0-50. the default value is 0 if not specified.
* `ttl` - (Optional, Int) Cache time. users can specify a value range of 60-86400. the smaller the value, the faster the modification records will take effect in all regions. default value: 300. unit: seconds.
* `weight` - (Optional, Int) DNS record weight. users can specify a value range of -1 to 100. a value of 0 means no resolution. if not specified, the default is -1, which means no weight is set. weight configuration is only applicable when type (dns record type) is A, AAAA, or CNAME. note: for the same subdomain, different dns records with the same resolution route should either all have weights set or none have weights set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `modified_on` - Modify time.
* `record_id` - DNS record ID.
* `status` - DNS record resolution status, the following values:
	- enable: has taken effect;
	- disable: has been disabled.


