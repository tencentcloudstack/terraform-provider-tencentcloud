---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_listener_default_domain"
sidebar_current: "docs-tencentcloud-resource-clb_listener_default_domain"
description: |-
  Provides a resource to set clb listener default domain
---

# tencentcloud_clb_listener_default_domain

Provides a resource to set clb listener default domain

## Example Usage

### Set default domain

```hcl
resource "tencentcloud_clb_listener_default_domain" "example" {
  clb_id      = "lb-g1miv1ok"
  listener_id = "lbl-duilx5qm"
  domain      = "3.com"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, String, ForceNew) ID of CLB instance.
* `domain` - (Required, String) Domain name of the listener rule. Single domain rules are passed to `domain`, and multi domain rules are passed to `domains`.
* `listener_id` - (Required, String, ForceNew) ID of CLB listener.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - ID of this CLB listener rule.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyDomainAttributes](https://www.tencentcloud.com/document/product/214/33819) - This API is used to modify domain-level attributes of Cloud Load Balancer layer-7 listener forwarding rules, such as modifying domain name, changing DefaultServer, enabling/disabling Http/2, and modifying certificates.
* Read: [DescribeListeners](https://www.tencentcloud.com/document/product/214/33831) - This API is used to get the list of listeners by CLB ID, listener protocol, or listener port. If no filter is specified, all listeners for the CLB instance will be returned.
* Update: [ModifyDomainAttributes](https://www.tencentcloud.com/document/product/214/33819) - This API is used to modify domain-level attributes of Cloud Load Balancer layer-7 listener forwarding rules, such as modifying domain name, changing DefaultServer, enabling/disabling Http/2, and modifying certificates.

## Import

CLB listener default domain can be imported using the id (version >= 1.47.0), e.g.

```
$ terraform import tencentcloud_clb_listener_default_domain.example lb-k2zjp9lv#lbl-hh141sn9
```

