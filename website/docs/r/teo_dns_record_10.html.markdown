---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record_10"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record_10"
description: |-
  Provides a resource to create a teo dns_record (version 10)
---

# tencentcloud_teo_dns_record_10

Provides a resource to create a teo dns_record (version 10)

## Example Usage

```hcl
# Basic DNS record creation
resource "tencentcloud_teo_dns_record_10" "example" {
  zone_id = "zone-xxxxxxxxxx"
  type    = "A"
  name    = "www.example.com"
  content = "192.168.1.1"
}

# DNS record with all optional parameters
resource "tencentcloud_teo_dns_record_10" "full_example" {
  zone_id  = "zone-xxxxxxxxxx"
  type     = "A"
  name     = "api.example.com"
  content  = "192.168.1.2"
  location = "Default"
  ttl      = 600
  weight   = 10
}

# MX record
resource "tencentcloud_teo_dns_record_10" "mx_example" {
  zone_id  = "zone-xxxxxxxxxx"
  type     = "MX"
  name     = "example.com"
  content  = "mail.example.com"
  priority = 10
}

# CNAME record
resource "tencentcloud_teo_dns_record_10" "cname_example" {
  zone_id  = "zone-xxxxxxxxxx"
  type     = "CNAME"
  name     = "www.example.com"
  content  = "example.com"
  location = "Default"
  ttl      = 300
}

# TXT record
resource "tencentcloud_teo_dns_record_10" "txt_example" {
  zone_id = "zone-xxxxxxxxxx"
  type    = "TXT"
  name    = "_dmarc.example.com"
  content = "v=DMARC1; p=none"
}

# AAAA record (IPv6)
resource "tencentcloud_teo_dns_record_10" "aaaa_example" {
  zone_id = "zone-xxxxxxxxxx"
  type    = "AAAA"
  name    = "ipv6.example.com"
  content = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
}

# DNS record with custom timeouts
resource "tencentcloud_teo_dns_record_10" "timeout_example" {
  zone_id = "zone-xxxxxxxxxx"
  type    = "A"
  name    = "slow.example.com"
  content = "192.168.1.3"

  timeouts {
    create = "15m"
    update = "10m"
    delete = "10m"
    read   = "5m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record content. Fill in corresponding content according to the type value. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `name` - (Required, String) DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `type` - (Required, String) DNS record type. Valid values are:
	- A: points domain name to an external IPv4 address, such as 8.8.8.8;
	- AAAA: points domain name to an external IPv6 address;
	- MX: used for email servers. When there are multiple MX records, lower priority value means higher priority;
	- CNAME: points domain name to another domain name, which then resolves to the final IP address;
	- TXT: identifies and describes domain name, commonly used for domain verification and SPF records (anti-spam);
	- NS: if you need to delegate subdomain to another DNS service provider for resolution, you need to add an NS record. Root domain cannot add NS records;
	- CAA: specifies CA that can issue certificates for this site;
	- SRV: identifies a server using a service, commonly used in Microsoft's directory management.
Different record types, such as SRV and CAA records, have different requirements for host record names and record value formats. For detailed descriptions and format examples of each record type, please refer to: [Introduction to DNS record types](https://cloud.tencent.com/document/product/1552/90453#2f681022-91ab-4a9e-ac3d-0a6c454d954e).
* `zone_id` - (Required, String, ForceNew) Zone ID.
* `location` - (Optional, String) DNS record resolution route. If not specified, default is Default, which means the default resolution route and takes effect in all regions.

- Resolution route configuration is only applicable when type (DNS record type) is A, AAAA, or CNAME.
- Resolution route configuration is only applicable to standard version and enterprise edition packages. For valid values, please refer to: [Resolution routes and corresponding code enumeration](https://cloud.tencent.com/document/product/1552/112542).
* `priority` - (Optional, Int) MX record priority, which takes effect only when type (DNS record type) is MX. Smaller value means higher priority. Users can specify a value range of 0-50. Default value is 0 if not specified.
* `ttl` - (Optional, Int) Cache time. Users can specify a value range of 60-86400. Smaller value means faster modification records will take effect in all regions. Default value: 300. Unit: seconds.
* `weight` - (Optional, Int) DNS record weight. Users can specify a value range of -1 to 100. A value of 0 means no resolution. If not specified, default is -1, which means no weight is set. Weight configuration is only applicable when type (DNS record type) is A, AAAA, or CNAME. Note: For the same subdomain, different DNS records with the same resolution route should either all have weights set or none have weights set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `modified_on` - Modification time.
* `record_id` - DNS record ID.
* `status` - DNS record resolution status. Valid values:
	- enable: has taken effect;
	- disable: has been disabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `10m`) Used when creating the resource.
* `read` - (Defaults to `5m`) Used when reading the resource.
* `update` - (Defaults to `10m`) Used when updating the resource.
* `delete` - (Defaults to `10m`) Used when deleting the resource.

## Import

teo_dns_record_10 can be imported using the zone_id and record_id, e.g.

```
terraform import tencentcloud_teo_dns_record_10.example zone-xxxxxxxxxx#record-xxxxxxxxxx
```

