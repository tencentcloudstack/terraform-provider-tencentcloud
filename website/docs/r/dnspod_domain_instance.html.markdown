---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_domain_instance"
sidebar_current: "docs-tencentcloud-resource-dnspod_domain_instance"
description: |-
  Provide a resource to create a DnsPod Domain instance.
---

# tencentcloud_dnspod_domain_instance

Provide a resource to create a DnsPod Domain instance.

## Example Usage

```hcl
resource "tencentcloud_dnspod_domain_instance" "foo" {
  domain = "hello.com"
  remark = "this is demo"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The Domain.
* `group_id` - (Optional, ForceNew) The Group Id of Domain.
* `is_mark` - (Optional, ForceNew) Whether to Mark the Domain.
* `remark` - (Optional) The remark of Domain.
* `status` - (Optional) The status of Domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the domain.


## Import

DnsPod Domain instance can be imported, e.g.

```
$ terraform import tencentcloud_dnspod_domain_instance.foo domain
```

