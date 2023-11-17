---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_domain_unlock"
sidebar_current: "docs-tencentcloud-resource-dnspod_domain_unlock"
description: |-
  Provides a resource to create a dnspod domain_unlock
---

# tencentcloud_dnspod_domain_unlock

Provides a resource to create a dnspod domain_unlock

## Example Usage

```hcl
resource "tencentcloud_dnspod_domain_unlock" "domain_unlock" {
  domain    = "dnspod.cn"
  lock_code = ""
  domain_id = 123
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `lock_code` - (Required, String, ForceNew) Domain unlock code, can be obtained through the ModifyDomainLock interface.
* `domain_id` - (Optional, Int, ForceNew) Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



