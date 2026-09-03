---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_records_status"
sidebar_current: "docs-tencentcloud-resource-teo_dns_records_status"
description: |-
  Provides a resource to manage TEO (EdgeOne) DNS records status config.
---

# tencentcloud_teo_dns_records_status

Provides a resource to manage TEO (EdgeOne) DNS records status config.

## Example Usage

### Enable a DNS record

```hcl
resource "tencentcloud_teo_dns_records_status" "example" {
  zone_id = "zone-3edjdliiw3he"

  records_to_enable = ["record-1234567890"]
}
```

### Disable a DNS record

```hcl
resource "tencentcloud_teo_dns_records_status" "example" {
  zone_id = "zone-3edjdliiw3he"

  records_to_disable = ["record-1234567890"]
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Site ID.
* `records_to_disable` - (Optional, List: [`String`]) DNS record ID list to be disabled, only manages a single resource, pass in a single record ID.
* `records_to_enable` - (Optional, List: [`String`]) DNS record ID list to be enabled, only manages a single resource, pass in a single record ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TEO (EdgeOne) DNS records status config can be imported using the `zone_id`, e.g.

```
terraform import tencentcloud_teo_dns_records_status.example zone-3edjdliiw3he
```

