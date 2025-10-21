---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_domain_alias"
sidebar_current: "docs-tencentcloud-resource-dnspod_domain_alias"
description: |-
  Provides a resource to create a dnspod domain_alias
---

# tencentcloud_dnspod_domain_alias

Provides a resource to create a dnspod domain_alias

## Example Usage

```hcl
resource "tencentcloud_dnspod_domain_alias" "domain_alias" {
  domain_alias = "dnspod.com"
  domain       = "dnspod.cn"
}
```

## Argument Reference

The following arguments are supported:

* `domain_alias` - (Required, String, ForceNew) Domain alias.
* `domain` - (Required, String, ForceNew) Domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `domain_alias_id` - Domain alias ID.


## Import

dnspod domain_alias can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_domain_alias.domain_alias domain#domain_alias_id
```

