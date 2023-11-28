---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_waf_domains"
sidebar_current: "docs-tencentcloud-resource-tse_waf_domains"
description: |-
  Provides a resource to create a tse waf_domains
---

# tencentcloud_tse_waf_domains

Provides a resource to create a tse waf_domains

## Example Usage

```hcl
resource "tencentcloud_tse_waf_domains" "waf_domains" {
  domain     = "tse.exmaple.com"
  gateway_id = "gateway-ed63e957"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) The waf protected domain name.
* `gateway_id` - (Required, String, ForceNew) Gateway ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tse waf_domains can be imported using the id, e.g.

```
terraform import tencentcloud_tse_waf_domains.waf_domains waf_domains_id
```

