---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_zone"
sidebar_current: "docs-tencentcloud-resource-private_dns_zone"
description: |-
  Provide a resource to create a Private Dns Zone.
---

# tencentcloud_private_dns_zone

Provide a resource to create a Private Dns Zone.

## Example Usage

```hcl
resource "tencentcloud_private_dns_zone" "foo" {
  domain = "domain.com"
  tags {
    "created_by" : "terraform"
  }
  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = "vpc-xxxxx"
  }
  remark             = "test"
  dns_forward_status = "DISABLED"
  account_vpc_set {
    uin         = "454xxxxxxx"
    region      = "ap-guangzhou"
    uniq_vpc_id = "vpc-xxxxx"
    vpc_name    = "test-redis"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) Domain name, which must be in the format of standard TLD.
* `account_vpc_set` - (Optional) List of authorized accounts' VPCs to associate with the private domain.
* `dns_forward_status` - (Optional) Whether to enable subdomain recursive DNS. Valid values: ENABLED, DISABLED. Default value: DISABLED.
* `remark` - (Optional) Remarks.
* `tag_set` - (Optional, **Deprecated**) It has been deprecated from version 1.72.4. Use `tags` instead. Tags the private domain when it is created.
* `tags` - (Optional) Tags of the private dns zone.
* `vpc_set` - (Optional) Associates the private domain to a VPC when it is created.

The `account_vpc_set` object supports the following:

* `region` - (Required) Region.
* `uin` - (Required) UIN of the VPC account.
* `uniq_vpc_id` - (Required) VPC ID.
* `vpc_name` - (Required) VPC NAME.

The `tag_set` object supports the following:

* `tag_key` - (Required) Key of Tag.
* `tag_value` - (Required) Value of Tag.

The `vpc_set` object supports the following:

* `region` - (Required) VPC REGION.
* `uniq_vpc_id` - (Required) VPC ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Private Dns Zone can be imported, e.g.

```
$ terraform import tencentcloud_private_dns_zone.foo zone_id
```

