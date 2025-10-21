---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_global_domain_dns"
sidebar_current: "docs-tencentcloud-resource-gaap_global_domain_dns"
description: |-
  Provides a resource to create a gaap global domain dns
---

# tencentcloud_gaap_global_domain_dns

Provides a resource to create a gaap global domain dns

## Example Usage

```hcl
resource "tencentcloud_gaap_global_domain_dns" "global_domain_dns" {
  domain_id                  = "dm-xxxxxx"
  proxy_id_list              = ["link-xxxxxx"]
  nation_country_inner_codes = ["101001"]
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, String) Domain Id.
* `nation_country_inner_codes` - (Required, Set: [`String`]) Nation Country Inner Codes.
* `proxy_id_list` - (Required, Set: [`String`]) Proxy Id List.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

gaap global_domain_dns can be imported using the id, e.g.

```
terraform import tencentcloud_gaap_global_domain_dns.global_domain_dns ${domainId}#${dnsRecordId}
```

