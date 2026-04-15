Provides a resource to create a teo dns_record (version 10)

Example Usage

```hcl
# Basic DNS record creation
resource "tencentcloud_teo_dns_record_10" "example" {
  zone_id  = "zone-xxxxxxxxxx"
  type     = "A"
  name     = "www.example.com"
  content  = "192.168.1.1"
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
  zone_id  = "zone-xxxxxxxxxx"
  type     = "TXT"
  name     = "_dmarc.example.com"
  content  = "v=DMARC1; p=none"
}

# AAAA record (IPv6)
resource "tencentcloud_teo_dns_record_10" "aaaa_example" {
  zone_id  = "zone-xxxxxxxxxx"
  type     = "AAAA"
  name     = "ipv6.example.com"
  content  = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
}

# DNS record with custom timeouts
resource "tencentcloud_teo_dns_record_10" "timeout_example" {
  zone_id  = "zone-xxxxxxxxxx"
  type     = "A"
  name     = "slow.example.com"
  content  = "192.168.1.3"

  timeouts {
    create = "15m"
    update = "10m"
    delete = "10m"
    read   = "5m"
  }
}
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Zone ID.
* `type` - (Required) DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.
* `name` - (Required) DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `content` - (Required) DNS record content. Fill in corresponding content according to the type value.
* `location` - (Optional) DNS record resolution route. Default: Default. Only applicable when type is A, AAAA, or CNAME.
* `ttl` - (Optional) Cache time in seconds. Range: 60-86400. Default: 300.
* `weight` - (Optional) DNS record weight. Range: -1-100. Default: -1. Only applicable when type is A, AAAA, or CNAME.
* `priority` - (Optional) MX record priority. Range: 0-50. Default: 0. Only applicable when type is MX.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `record_id` - DNS record ID.
* `status` - DNS record resolution status. Valid values: enable, disable.
* `created_on` - Creation time.
* `modified_on` - Modification time.

Import

teo_dns_record_10 can be imported using the zone_id and record_id, e.g.

```
terraform import tencentcloud_teo_dns_record_10.example zone-xxxxxxxxxx#record-xxxxxxxxxx
```
