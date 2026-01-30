---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_line_group"
sidebar_current: "docs-tencentcloud-resource-dnspod_line_group"
description: |-
  Provides a resource to create a DNSPod line group.
---

# tencentcloud_dnspod_line_group

Provides a resource to create a DNSPod line group.

## Example Usage

```hcl
resource "tencentcloud_dnspod_line_group" "example" {
  domain = "example.com"
  name   = "telecom_group"
  lines  = ["电信", "移动"]
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `lines` - (Required, List: [`String`]) List of lines in the group. Maximum 120 lines.
* `name` - (Required, String) Line group name, length 1-17 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Creation time.
* `domain_id` - Domain ID.
* `line_group_id` - Line group ID.
* `updated_on` - Update time.


## Import

DNSPod line group can be imported using the id (format: `{domain}#{line_group_id}`), e.g.

```
$ terraform import tencentcloud_dnspod_line_group.example example.com#123
```

