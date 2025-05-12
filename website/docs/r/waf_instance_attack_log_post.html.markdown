---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_instance_attack_log_post"
sidebar_current: "docs-tencentcloud-resource-waf_instance_attack_log_post"
description: |-
  Provides a resource to create a WAF instance attack log post
---

# tencentcloud_waf_instance_attack_log_post

Provides a resource to create a WAF instance attack log post

~> **NOTE:** Only enterprise version and above are supported for activation

## Example Usage

```hcl
resource "tencentcloud_waf_instance_attack_log_post" "example" {
  instance_id     = "waf_2kxtlbky11b4wcrb"
  attack_log_post = 1
}
```

## Argument Reference

The following arguments are supported:

* `attack_log_post` - (Required, Int) Attack log delivery switch. 0- Disable, 1- Enable.
* `instance_id` - (Required, String, ForceNew) Waf instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

WAF instance attack log post can be imported using the id, e.g.

```
terraform import tencentcloud_waf_instance_attack_log_post.example waf_2kxtlbky11b4wcrb
```

