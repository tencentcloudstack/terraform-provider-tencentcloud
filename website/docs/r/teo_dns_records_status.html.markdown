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

  filters {
    name   = "id"
    values = ["record-1234567890"]
  }

  records_to_enable = ["record-1234567890"]
}
```

### Disable a DNS record

```hcl
resource "tencentcloud_teo_dns_records_status" "example" {
  zone_id = "zone-3edjdliiw3he"

  filters {
    name   = "id"
    values = ["record-1234567890"]
  }

  records_to_disable = ["record-1234567890"]
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Site ID.
* `filters` - (Optional, List) Filter conditions, each filter element contains `name`, `values` and `fuzzy`.
* `match` - (Optional, String) Match mode, valid values: `all`, `any`.
* `records_to_disable` - (Optional, List: [`String`]) DNS record ID list to be disabled, only manages a single resource, pass in a single record ID.
* `records_to_enable` - (Optional, List: [`String`]) DNS record ID list to be enabled, only manages a single resource, pass in a single record ID.
* `sort_by` - (Optional, String) Sort by field.
* `sort_order` - (Optional, String) Sort order, valid values: `asc`, `desc`.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, List) Filter values of the field.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `dns_records` - DNS record list.
  * `content` - DNS record content.
  * `created_on` - Creation time.
  * `location` - DNS record resolution line.
  * `modified_on` - Modification time.
  * `name` - DNS record name.
  * `priority` - MX record priority.
  * `record_id` - DNS record ID.
  * `status` - DNS record resolution status, valid values: `enable`, `disable`.
  * `ttl` - Cache time, unit: seconds.
  * `type` - DNS record type.
  * `weight` - DNS record weight.
  * `zone_id` - Site ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `read` - (Defaults to `10m`) Used when reading the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `10m`) Used when deleting the resource.

## Import

TEO (EdgeOne) DNS records status config can be imported using the `zone_id`, e.g.

```
terraform import tencentcloud_teo_dns_records_status.example zone-3edjdliiw3he
```

