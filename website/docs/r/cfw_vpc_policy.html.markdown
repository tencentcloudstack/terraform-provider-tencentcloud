---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_vpc_policy"
sidebar_current: "docs-tencentcloud-resource-cfw_vpc_policy"
description: |-
  Provides a resource to create a cfw vpc_policy
---

# tencentcloud_cfw_vpc_policy

Provides a resource to create a cfw vpc_policy

## Example Usage

```hcl
resource "tencentcloud_cfw_vpc_policy" "example" {
  source_content = "0.0.0.0/0"
  source_type    = "net"
  dest_content   = "192.168.0.2"
  dest_type      = "net"
  protocol       = "ANY"
  rule_action    = "log"
  port           = "-1/-1"
  description    = "description."
  enable         = "true"
  fw_group_id    = "ALL"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required, String) Describe.
* `dest_content` - (Required, String) Access purpose example: net:IP/CIDR(192.168.0.2) domain:domain rule, for example*.qq.com.
* `dest_type` - (Required, String) Access purpose type, the type can be: net, template.
* `port` - (Required, String) The port for the access control policy. Value: -1/-1: All ports; 80: port 80.
* `protocol` - (Required, String) Protocol, optional value:TCP, UDP, ICMP, ANY, HTTP, HTTPS, HTTP/HTTPS, SMTP, SMTPS, SMTP/SMTPS, FTP, DNS, TLS/SSL.
* `rule_action` - (Required, String) How traffic set in the access control policy passes through the cloud firewall. Value: accept:accept, drop:drop, log:log.
* `source_content` - (Required, String) Access source examplnet:IP/CIDR(192.168.0.2).
* `source_type` - (Required, String) Access source type, the type can be: net, template.
* `enable` - (Optional, String) Rule status, true means enabled, false means disabled. Default is true.
* `fw_group_id` - (Optional, String) Firewall instance ID where the rule takes effect. Default is ALL.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `fw_group_name` - Firewall name.
* `internal_uuid` - Uuid used internally, this field is generally not used.
* `uuid` - The unique id corresponding to the rule.


## Import

cfw vpc_policy can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_vpc_policy.vpc_policy vpc_policy_id
```

