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

### Create a basic Private Dns Zone

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_private_dns_zone" "example" {
  domain = "domain.com"
  remark = "remark."

  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = tencentcloud_vpc.vpc.id
  }

  dns_forward_status   = "DISABLED"
  cname_speedup_status = "ENABLED"

  tags = {
    createdBy : "terraform"
  }
}
```

### Create a Private Dns Zone domain and bind associated accounts'VPC

```hcl
resource "tencentcloud_private_dns_zone" "example" {
  domain = "domain.com"
  remark = "remark."

  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = tencentcloud_vpc.vpc.id
  }

  account_vpc_set {
    uin         = "123456789"
    uniq_vpc_id = "vpc-adsebmya"
    region      = "ap-guangzhou"
    vpc_name    = "vpc-name"
  }

  dns_forward_status   = "DISABLED"
  cname_speedup_status = "ENABLED"

  tags = {
    createdBy : "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain name, which must be in the format of standard TLD.
* `account_vpc_set` - (Optional, List) List of authorized accounts' VPCs to associate with the private domain.
* `cname_speedup_status` - (Optional, String) CNAME acceleration: ENABLED, DISABLED, Default value is ENABLED.
* `dns_forward_status` - (Optional, String) Whether to enable subdomain recursive DNS. Valid values: ENABLED, DISABLED. Default value: DISABLED.
* `remark` - (Optional, String) Remarks.
* `tag_set` - (Optional, List, **Deprecated**) It has been deprecated from version 1.72.4. Use `tags` instead. Tags the private domain when it is created.
* `tags` - (Optional, Map) Tags of the private dns zone.
* `vpc_set` - (Optional, List) Associates the private domain to a VPC when it is created.

The `account_vpc_set` object supports the following:

* `region` - (Required, String) Region.
* `uin` - (Required, String) UIN of the VPC account.
* `uniq_vpc_id` - (Required, String) VPC ID.
* `vpc_name` - (Required, String) VPC NAME.

The `tag_set` object supports the following:

* `tag_key` - (Required, String) Key of Tag.
* `tag_value` - (Required, String) Value of Tag.

The `vpc_set` object supports the following:

* `region` - (Required, String) VPC REGION.
* `uniq_vpc_id` - (Required, String) VPC ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Private Dns Zone can be imported, e.g.

```
$ terraform import tencentcloud_private_dns_zone.foo zone_id
```

