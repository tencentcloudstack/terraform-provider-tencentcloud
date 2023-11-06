---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_modify_record_group_operation"
sidebar_current: "docs-tencentcloud-resource-dnspod_modify_record_group_operation"
description: |-
  Provides a resource to create a dnspod tencentcloud_dnspod_modify_record_group_operation
---

# tencentcloud_dnspod_modify_record_group_operation

Provides a resource to create a dnspod tencentcloud_dnspod_modify_record_group_operation

## Example Usage

```hcl
resource "tencentcloud_dnspod_modify_record_group_operation" "modify_record_group" {
  domain    = "dnspod.cn"
  group_id  = 1
  record_id = "234|345"
  domain_id = 123
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain.
* `group_id` - (Required, Int, ForceNew) Record Group ID.
* `record_id` - (Required, String, ForceNew) Record ID, multiple IDs are separated by a vertical line |.
* `domain_id` - (Optional, Int, ForceNew) Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



