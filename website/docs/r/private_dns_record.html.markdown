---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_record"
sidebar_current: "docs-tencentcloud-resource-private_dns_record"
description: |-
  Provide a resource to create a Private Dns Record.
---

# tencentcloud_private_dns_record

Provide a resource to create a Private Dns Record.

## Example Usage

```hcl
# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

# create private dns zone
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
    createdBy = "Terraform"
  }
}

# create private dns record
resource "tencentcloud_private_dns_record" "example" {
  zone_id      = tencentcloud_private_dns_zone.example.id
  record_type  = "A"
  record_value = "192.168.1.2"
  sub_domain   = "www"
  ttl          = 300
  weight       = 20
  mx           = 0
  status       = "disabled"
}
```

## Argument Reference

The following arguments are supported:

* `record_type` - (Required, String) Record type. Valid values: `A`, `AAAA`, `CNAME`, `MX`, `TXT`, `PTR`.
* `record_value` - (Required, String) Record value, such as IP: 192.168.10.2, CNAME: cname.qcloud.com, and MX: mail.qcloud.com.
* `sub_domain` - (Required, String) Subdomain, such as `www`, `m`, and `@`.
* `zone_id` - (Required, String, ForceNew) Private domain ID.
* `mx` - (Optional, Int) MX priority, which is required when the record type is MX. Valid values: 5, 10, 15, 20, 30, 40, 50.
* `status` - (Optional, String) Record status. Valid values: `enabled`, `disabled`.
* `ttl` - (Optional, Int) Record cache time. The smaller the value, the faster the record will take effect. Value range: 1~86400s.
* `weight` - (Optional, Int) Record weight. Value range: 1~100.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Private Dns Record can be imported, e.g.

```
$ terraform import tencentcloud_private_dns_record.example zone-iza3a33s#1983030
```

