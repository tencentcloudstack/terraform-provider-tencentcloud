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
resource "tencentcloud_private_dns_record" "foo" {
  zone_id      = "zone-rqndjnki"
  record_type  = "A"
  record_value = "192.168.1.2"
  sub_domain   = "www"
  ttl          = 300
  weight       = 1
  mx           = 0
}
```

## Argument Reference

The following arguments are supported:

* `record_type` - (Required) Record type. Valid values: "A", "AAAA", "CNAME", "MX", "TXT", "PTR".
* `record_value` - (Required) Record value, such as IP: 192.168.10.2, CNAME: cname.qcloud.com, and MX: mail.qcloud.com..
* `sub_domain` - (Required) Subdomain, such as "www", "m", and "@".
* `zone_id` - (Required) Private domain ID.
* `mx` - (Optional) MX priority, which is required when the record type is MX. Valid values: 5, 10, 15, 20, 30, 40, 50.
* `ttl` - (Optional) Record cache time. The smaller the value, the faster the record will take effect. Value range: 1~86400s.
* `weight` - (Optional) Record weight. Value range: 1~100.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Private Dns Record can be imported, e.g.

```
$ terraform import tencentcloud_private_dns_zone.foo zone_id#record_id
```

