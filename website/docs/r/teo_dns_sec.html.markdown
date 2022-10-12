---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_sec"
sidebar_current: "docs-tencentcloud-resource-teo_dns_sec"
description: |-
  Provides a resource to create a teo dns_sec
---

# tencentcloud_teo_dns_sec

Provides a resource to create a teo dns_sec

## Example Usage

```hcl
resource "tencentcloud_teo_dns_sec" "dns_sec" {
  zone_id = "zone-297z8rf93cfw"
  status  = "enabled"
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Required, String) DNSSEC status. Valid values: `enabled`, `disabled`.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `dnssec` - (Optional, List) DNSSEC infos.

The `dnssec` object supports the following:


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `modified_on` - Last modification date.


## Import

teo dns_sec can be imported using the zone_id, e.g.
```
$ terraform import tencentcloud_teo_dns_sec.dns_sec zone-297z8rf93cfw
```

