---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_modify_domain_owner_operation"
sidebar_current: "docs-tencentcloud-resource-dnspod_modify_domain_owner_operation"
description: |-
  Provides a resource to create a dnspod modify_domain_owner
---

# tencentcloud_dnspod_modify_domain_owner_operation

Provides a resource to create a dnspod modify_domain_owner

## Example Usage

```hcl
resource "tencentcloud_dnspod_modify_domain_owner_operation" "modify_domain_owner" {
  domain    = "dnspod.cn"
  account   = "xxxxxxxxx"
  domain_id = 123
}
```

## Argument Reference

The following arguments are supported:

* `account` - (Required, String, ForceNew) The account to which the domain needs to be transferred, supporting Uin or email format.
* `domain` - (Required, String, ForceNew) Domain.
* `domain_id` - (Optional, Int, ForceNew) Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



