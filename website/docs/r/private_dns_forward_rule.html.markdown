---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_forward_rule"
sidebar_current: "docs-tencentcloud-resource-private_dns_forward_rule"
description: |-
  Provides a resource to create a privatedns forward rule
---

# tencentcloud_private_dns_forward_rule

Provides a resource to create a privatedns forward rule

## Example Usage

```hcl
resource "tencentcloud_private_dns_end_point" "example" {
  end_point_name       = "tf-example"
  end_point_service_id = "vpcsvc-61wcwmar"
  end_point_region     = "ap-guangzhou"
  ip_num               = 1
}

resource "tencentcloud_private_dns_forward_rule" "example" {
  rule_name    = "tf-example"
  rule_type    = "DOWN"
  zone_id      = "zone-cmmbvaq8"
  end_point_id = tencentcloud_private_dns_end_point.example.id
}
```

## Argument Reference

The following arguments are supported:

* `end_point_id` - (Required, String) Endpoint ID.
* `rule_name` - (Required, String) Forwarding rule name.
* `rule_type` - (Required, String) Forwarding rule type. DOWN: From cloud to off-cloud; UP: From off-cloud to cloud.
* `zone_id` - (Required, String) Private domain ID, which can be viewed on the private domain list page.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

private dns forward rule can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_forward_rule.example fid-dbc2c0a97c
```

