---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_scheduling_domain_user_name"
sidebar_current: "docs-tencentcloud-resource-antiddos_scheduling_domain_user_name"
description: |-
  Provides a resource to create a antiddos scheduling_domain_user_name
---

# tencentcloud_antiddos_scheduling_domain_user_name

Provides a resource to create a antiddos scheduling_domain_user_name

## Example Usage

```hcl
resource "tencentcloud_antiddos_scheduling_domain_user_name" "scheduling_domain_user_name" {
  domain_name      = "test.com"
  domain_user_name = ""
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String, ForceNew) user cname.
* `domain_user_name` - (Required, String) domain name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

antiddos scheduling_domain_user_name can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_scheduling_domain_user_name.scheduling_domain_user_name ${domainName}
```

