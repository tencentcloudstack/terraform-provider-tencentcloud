---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record"
description: |-
  Provides a resource to create a teo teo_dns_record
---

# tencentcloud_teo_dns_record

Provides a resource to create a teo teo_dns_record

## Example Usage

```hcl
resource "tencentcloud_teo_dns_record" "teo_dns_record" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  name     = "a.makn.cn"
  priority = 5
  ttl      = 300
  weight   = -1
  status   = "enable"
}
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
* `status` - (Optional, String) DNS record resolution status, the following values:
	- enable: has taken effect;
	- disable: has been disabled.
* `ttl` - (Optional, Int) Cache time. users can specify a value range of 60-86400. the smaller the value, the faster the modification records will take effect in all regions. default value: 300. unit: seconds.
* `weight` - (Optional, Int) DNS record weight. users can specify a value range of -1 to 100. a value of 0 means no resolution. if not specified, the default is -1, which means no weight is set. weight configuration is only applicable when type (dns record type) is A, AAAA, or CNAME. note: for the same subdomain, different dns records with the same resolution route should either all have weights set or none have weights set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `modified_on` - Modify time.


## Import

teo teo_dns_record can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record.teo_dns_record {zoneId}#{recordId}
```

