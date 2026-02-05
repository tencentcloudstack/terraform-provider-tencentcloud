---
subcategory: "DNSPod"
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

# Access computed fields
output "domain_status" {
  value = tencentcloud_dnspod_domain_instance.foo.status
}

output "record_count" {
  value = tencentcloud_dnspod_domain_instance.foo.record_count
}

output "domain_grade" {
  value = tencentcloud_dnspod_domain_instance.foo.grade
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) The Domain.
* `group_id` - (Optional, Int, ForceNew) The Group Id of Domain.
* `is_mark` - (Optional, String, ForceNew) Whether to Mark the Domain.
* `remark` - (Optional, String) The remark of Domain.
* `status` - (Optional, String) The status of Domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the domain.
* `domain_id` - ID of the domain.
* `grade` - The DNS plan/package grade of the domain (e.g., DP_Free, DP_Plus).
* `record_count` - Number of DNS records under this domain.
* `slave_dns` - Is secondary DNS enabled.
* `updated_on` - Last modification time of the domain.


## Import

DnsPod Domain instance can be imported, e.g.

```
$ terraform import tencentcloud_dnspod_domain_instance.foo domain
```

