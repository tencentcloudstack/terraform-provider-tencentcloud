---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_edge_policy"
sidebar_current: "docs-tencentcloud-resource-cfw_edge_policy"
description: |-
  Provides a resource to create a CFW edge policy
---

# tencentcloud_cfw_edge_policy

Provides a resource to create a CFW edge policy

## Example Usage

```hcl
resource "tencentcloud_cfw_edge_policy" "example" {
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
  scope          = "all"
}
```

### If target_type is tag

```hcl
resource "tencentcloud_cfw_edge_policy" "example" {
  source_content = "0.0.0.0/0"
  source_type    = "net"
  target_content = jsonencode({ "Key" : "test", "Value" : "dddd" })
  target_type    = "tag"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, Int) Rule direction: 1, inbound; 0, outbound.
* `port` - (Required, String) The port for the access control policy. Value: -1/-1: All ports 80: Port 80.
* `protocol` - (Required, String) Protocol. If Direction=1 && Scope=serial, optional values: TCP UDP ICMP ANY HTTP HTTPS HTTP/HTTPS SMTP SMTPS SMTP/SMTPS FTP DNS; If Direction=1 && Scope!=serial, optional values: TCP; If Direction=0 && Scope=serial, optional values: TCP UDP ICMP ANY HTTP HTTPS HTTP/HTTPS SMTP SMTPS SMTP/SMTPS FTP DNS; If Direction=0 && Scope!=serial, optional values: TCP HTTP/HTTPS TLS/SSL.
* `rule_action` - (Required, String) How the traffic set in the access control policy passes through the cloud firewall. Values: accept: allow; drop: reject; log: observe.
* `source_content` - (Required, String) Access source example: net:IP/CIDR(192.168.0.2).
* `source_type` - (Required, String) Access source type: for inbound rules, the type can be net, location, vendor, template; for outbound rules, it can be net, instance, tag, template, group.
* `target_content` - (Required, String) Example of access purpose: net: IP/CIDR(192.168.0.2) domain: domain name rules, such as *.qq.com.
* `target_type` - (Required, String) Access purpose type: For inbound rules, the type can be net, instance, tag, template, group; for outbound rules, it can be net, location, vendor, template.
* `description` - (Optional, String) Description.
* `enable` - (Optional, String) Rule status, true means enabled, false means disabled. Default is true.
* `param_template_id` - (Optional, String) Parameter template id.
* `scope` - (Optional, String) Effective range. serial: serial; side: bypass; all: global, Default is all.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `uuid` - The unique id corresponding to the rule, no need to fill in when creating the rule.


## Import

CFW edge policy can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_edge_policy.example 1859582
```

