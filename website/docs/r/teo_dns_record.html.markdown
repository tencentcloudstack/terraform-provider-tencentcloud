---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_record"
sidebar_current: "docs-tencentcloud-resource-teo_dns_record"
description: |-
  Provides a resource to create a teo dnsRecord
---

# tencentcloud_teo_dns_record

Provides a resource to create a teo dnsRecord

## Example Usage

```hcl
resource "tencentcloud_teo_dns_record" "dns_record" {
  zone_id     = tencentcloud_teo_zone.zone.id
  record_type = "A"
  name        = "sfurnace.work"
  mode        = "proxied"
  content     = "2.2.2.2"
  ttl         = 80
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) DNS Record Content.
* `mode` - (Required, String) Proxy mode. Valid values: dns_only, cdn_only, and secure_cdn.
* `name` - (Required, String) DNS Record Name.
* `record_type` - (Required, String) DNS Record Type.
* `zone_id` - (Required, String) Site ID.
* `priority` - (Optional, Int) Priority.
* `tags` - (Optional, Map) Tag description list.
* `ttl` - (Optional, Int) TTL, the range is 1-604800, and the minimum value of different levels of domain names is different.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname` - CNAME address.
* `created_on` - Creation time.
* `domain_status` - .
* `locked` - Whether the DNS record is locked.
* `modified_on` - Modification time.
* `status` - Resolution status.
* `zone_name` - Site Name.


## Import

teo dns_record can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_dns_record.dnsRecord zoneId#dnsRecordId#name
```

