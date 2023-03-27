---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record"
description: |-
  Provides a resource to create a teo dns_record
---

# tencentcloud_teo_dns_record

Provides a resource to create a teo dns_record

~> **NOTE:** This resource has been deprecated in Terraform TencentCloud Provider Version 1.79.19.

## Example Usage

```hcl
resource "tencentcloud_teo_dns_record" "dns_record" {
  zone_id  = "zone-297z8rf93cfw"
  type     = "A"
  name     = "www.toutiao2.com"
  content  = "150.109.8.2"
  mode     = "proxied"
  ttl      = "1"
  priority = 1
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS record Content.
* `mode` - (Required, String) Proxy mode. Valid values:- `dns_only`: only DNS resolution of the subdomain is enabled.- `proxied`: subdomain is proxied and accelerated.
* `name` - (Required, String) DNS record Name.
* `type` - (Required, String) DNS record Type. Valid values: `A`, `AAAA`, `CNAME`, `MX`, `TXT`, `NS`, `CAA`, `SRV`.
* `zone_id` - (Required, String) Site ID.
* `priority` - (Optional, Int) Priority of the record. Valid value range: 1-50, the smaller value, the higher priority.
* `status` - (Optional, String) Resolution status. Valid values: `active`, `pending`.
* `ttl` - (Optional, Int) Time to live of the DNS record cache in seconds.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname` - CNAME address. Note: This field may return null, indicating that no valid value can be obtained.
* `created_on` - Creation date.
* `dns_record_id` - DNS record ID.
* `domain_status` - Whether this domain enable load balancing, security, or l4 proxy capability. Valid values: `lb`, `security`, `l4`.
* `locked` - Whether the DNS record is locked.
* `modified_on` - Last modification date.


## Import

teo dns_record can be imported using the zone_id#dns_record_id, e.g.
```
$ terraform import tencentcloud_teo_dns_record.dns_record zone-297z8rf93cfw#record-297z9ei9b9oc
```

