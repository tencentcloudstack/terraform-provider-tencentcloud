---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_sec"
sidebar_current: "docs-tencentcloud-resource-teo_dns_sec"
description: |-
  Provides a resource to create a teo dnsSec
---

# tencentcloud_teo_dns_sec

Provides a resource to create a teo dnsSec

## Example Usage

```hcl
resource "tencentcloud_teo_dns_sec" "dns_sec" {
  zone_id = tencentcloud_teo_zone.zone.id
  status  = "disabled"
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Required, String) DNSSEC status. Valid values: `enabled`, `disabled`.
* `zone_id` - (Required, String) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `dnssec` - DNSSEC infos.
  * `algorithm` - Encryption algorithm.
  * `d_s` - DS record value.
  * `digest_algorithm` - Digest algorithm.
  * `digest_type` - Digest type.
  * `digest` - Digest message.
  * `flags` - Flag.
  * `key_tag` - Key tag.
  * `key_type` - Encryption type.
  * `public_key` - Public key.
* `modified_on` - Last modification date.
* `zone_name` - Site Name.


## Import

teo dns_sec can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_dns_sec.dns_sec zoneId
```

