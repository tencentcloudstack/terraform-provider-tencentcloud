---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_waf_protection"
sidebar_current: "docs-tencentcloud-resource-tse_waf_protection"
description: |-
  Provides a resource to create a tse waf_protection
---

# tencentcloud_tse_waf_protection

Provides a resource to create a tse waf_protection

## Example Usage

```hcl
resource "tencentcloud_tse_waf_protection" "waf_protection" {
  gateway_id = "gateway-ed63e957"
  type       = "Route"
  list       = ["7324a769-9d87-48ce-a904-48c3defc4abd"]
  operate    = "open"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Gateway ID.
* `operate` - (Required, String) `open`: open the protection, `close`: close the protection.
* `type` - (Required, String) The type of protection resource. Reference value: `Global`: instance, `Service`: service, `Route`: route, `Object`: obejct (This interface does not currently support this type).
* `list` - (Optional, Set: [`String`]) Means the list of services or routes when the resource type `Type` is `Service` or `Route`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `global_status` - Global protection status.


