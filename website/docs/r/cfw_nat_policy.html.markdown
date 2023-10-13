---
subcategory: "Cfw"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_nat_policy"
sidebar_current: "docs-tencentcloud-resource-cfw_nat_policy"
description: |-
  Provides a resource to create a cfw nat_policy
---

# tencentcloud_cfw_nat_policy

Provides a resource to create a cfw nat_policy

## Example Usage

```hcl
resource "tencentcloud_cfw_nat_policy" "example" {
  source_content = "1.1.1.1/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "policy description."
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, Int) Rule direction: 1, inbound; 0, outbound.
* `port` - (Required, String) The port for the access control policy. Value: -1/-1: All ports 80: Port 80.
* `protocol` - (Required, String) Protocol. If Direction=1, optional values: TCP, UDP, ANY; If Direction=0, optional values: TCP, UDP, ICMP, ANY, HTTP, HTTPS, HTTP/HTTPS, SMTP, SMTPS, SMTP/SMTPS, FTP, and DNS.
* `rule_action` - (Required, String) How the traffic set in the access control policy passes through the cloud firewall. Values: accept: allow; drop: reject; log: observe.
* `source_content` - (Required, String) Access source example: net:IP/CIDR(192.168.0.2).
* `source_type` - (Required, String) Access source type: for inbound rules, the type can be net, location, vendor, template; for outbound rules, it can be net, instance, tag, template, group.
* `target_content` - (Required, String) Example of access purpose: net: IP/CIDR(192.168.0.2) domain: domain name rules, such as *.qq.com.
* `target_type` - (Required, String) Access purpose type: For inbound rules, the type can be net, instance, tag, template, group; for outbound rules, it can be net, location, vendor, template.
* `description` - (Optional, String) Description.
* `enable` - (Optional, String) Rule status, true means enabled, false means disabled. Default is true.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `uuid` - The unique id corresponding to the rule, no need to fill in when creating the rule.


## Import

cfw nat_policy can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_nat_policy.example nat_policy_id
```

