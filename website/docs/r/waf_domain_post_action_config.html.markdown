---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_domain_post_action_config"
sidebar_current: "docs-tencentcloud-resource-waf_domain_post_action_config"
description: |-
  Provides a resource to create a WAF domain post action config
---

# tencentcloud_waf_domain_post_action_config

Provides a resource to create a WAF domain post action config

## Example Usage

```hcl
resource "tencentcloud_waf_domain_post_action_config" "example" {
  domain             = "example.com"
  post_cls_action    = 1
  post_ckafka_action = 0
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain.
* `post_ckafka_action` - (Required, Int) 0- Disable shipping, 1- Enable shipping.
* `post_cls_action` - (Required, Int) 0- Disable shipping, 1- Enable shipping.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

WAF domain post action config can be imported using the id, e.g.

```
terraform import tencentcloud_waf_domain_post_action_config.example example.com
```

