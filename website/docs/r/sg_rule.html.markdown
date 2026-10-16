---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sg_rule"
sidebar_current: "docs-tencentcloud-resource-sg_rule"
description: |-
  Provides a resource to create a cfw sg_rule
---

# tencentcloud_sg_rule

Provides a resource to create a cfw sg_rule

## Example Usage

```hcl
resource "tencentcloud_sg_rule" "sg_rule" {
  enable = 1

  data {
    description         = "1111112"
    dest_content        = "0.0.0.0/0"
    dest_type           = "net"
    port                = "-1/-1"
    protocol            = "ANY"
    rule_action         = "accept"
    service_template_id = "ppm-l9u5pf1y"
    source_content      = "0.0.0.0/0"
    source_type         = "net"
  }
}
```

## Argument Reference

The following arguments are supported:

* `data` - (Required, List) Creates rule data.
* `enable` - (Optional, Int) Rule status. `0` is off, `1` is on. This parameter is not required or is 1 when creating.

The `data` object supports the following:

* `description` - (Required, String) Description.
* `dest_content` - (Required, String) Destination example: `net`: IP/CIDR (192.168.0.2); `template`: parameter template (ipm-dyodhpby); `instance`: asset instance (ins-123456); `resourcegroup`: asset group (/all groups/group 1/subgroup 1); `tag`: resource tag ({"Key":"tag key","Value":"tag value"}); `region`: region (ap-gaungzhou).
* `dest_type` - (Required, String) Access destination type. Valid values: net|template|instance|resourcegroup|tag|region.
* `rule_action` - (Required, String) The action that Cloud Firewall performs on the traffic. Valid values: `accept`: allow, `drop`: deny.
* `source_content` - (Required, String) Source example: `net`: IP/CIDR (192.168.0.2); `template`: parameter template (ipm-dyodhpby); `instance`: asset instance (ins-123456); `resourcegroup`: asset group (/all groups/group 1/subgroup 1); `tag`: resource tag ({"Key":"tag key","Value":"tag value"}); `region`: region (ap-gaungzhou).
* `source_type` - (Required, String) Access source type. Valid values: net|template|instance|resourcegroup|tag|region.
* `port` - (Optional, String) The port to apply access control rules. Valid values: `-1/-1`: all ports, `80`: port 80.
* `protocol` - (Optional, String) Protocol. TCP/UDP/ICMP/ANY.
* `service_template_id` - (Optional, String) Parameter template ID of port and protocol type; mutually exclusive with Protocol and Port.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.

The `data` object exports the following:

* `order_index` - Rule priority.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddEnterpriseSecurityGroupRules](https://cloud.tencent.com/document/api/1132/64234)
* Create: [DescribeEnterpriseSGRuleProgress](https://cloud.tencent.com/document/api/1132/93416)
* Read: [DescribeEnterpriseSecurityGroupRuleList](https://cloud.tencent.com/document/api/1132/110987)
* Update: [DescribeEnterpriseSGRuleProgress](https://cloud.tencent.com/document/api/1132/93416)
* Update: [ModifyEnterpriseSecurityGroupRule](https://cloud.tencent.com/document/api/1132/86971)
* Delete: [DescribeEnterpriseSGRuleProgress](https://cloud.tencent.com/document/api/1132/93416)
* Delete: [RemoveEnterpriseSecurityGroupRule](https://cloud.tencent.com/document/api/1132/64232)

## Import

cfw sg_rule can be imported using the id, e.g.

```
terraform import tencentcloud_sg_rule.sg_rule rule_id
```

