---
subcategory: "DNSPod"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_package_domain"
sidebar_current: "docs-tencentcloud-resource-dnspod_package_domain"
description: |-
  Provides a resource to manage DNSPod package domain binding
---

# tencentcloud_dnspod_package_domain

Provides a resource to manage DNSPod package domain binding

## Example Usage

```hcl
resource "tencentcloud_dnspod_package_domain" "example" {
  resource_id = "91d8006a"
  domain_id   = 92435817
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int) Domain ID to bind to the package.
* `resource_id` - (Required, String, ForceNew) Package resource ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `domain` - Domain.


## Import

dnspod package_domain can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_package_domain.example 91d8006a
```

