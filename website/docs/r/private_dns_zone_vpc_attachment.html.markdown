---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_zone_vpc_attachment"
sidebar_current: "docs-tencentcloud-resource-private_dns_zone_vpc_attachment"
description: |-
  Provides a resource to create a PrivateDns zone_vpc_attachment
---

# tencentcloud_private_dns_zone_vpc_attachment

Provides a resource to create a PrivateDns zone_vpc_attachment

~> **NOTE:**  If you need to bind account A to account B's VPC resources, you need to first grant role authorization to account A.

## Example Usage

### Append VPC associated with private dns zone

```hcl
resource "tencentcloud_private_dns_zone" "example" {
  domain = "domain.com"
  remark = "remark."

  dns_forward_status   = "DISABLED"
  cname_speedup_status = "ENABLED"

  tags = {
    createdBy : "terraform"
  }
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_private_dns_zone_vpc_attachment" "example" {
  zone_id = tencentcloud_private_dns_zone.example.id

  vpc_set {
    uniq_vpc_id = tencentcloud_vpc.vpc.id
    region      = "ap-guangzhou"
  }
}
```

### Add VPC information for associated accounts in the private dns zone

```hcl
resource "tencentcloud_private_dns_zone_vpc_attachment" "example" {
  zone_id = tencentcloud_private_dns_zone.example.id

  account_vpc_set {
    uniq_vpc_id = "vpc-82znjzn3"
    region      = "ap-guangzhou"
    uin         = "100017155920"
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) PrivateZone ID.
* `account_vpc_set` - (Optional, List, ForceNew) New add account vpc info.
* `vpc_set` - (Optional, List, ForceNew) New add vpc info.

The `account_vpc_set` object supports the following:

* `region` - (Required, String) Vpc region.
* `uin` - (Required, String) Vpc owner uin. To grant role authorization to this account.
* `uniq_vpc_id` - (Required, String) Uniq Vpc Id.

The `vpc_set` object supports the following:

* `region` - (Required, String) Vpc region.
* `uniq_vpc_id` - (Required, String) Uniq Vpc Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

PrivateDns zone_vpc_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_zone_vpc_attachment.example zone-6t11lof0#vpc-jdx11z0t
```

