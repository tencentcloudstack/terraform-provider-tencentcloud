---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_cngw_certificate"
sidebar_current: "docs-tencentcloud-resource-tse_cngw_certificate"
description: |-
  Provides a resource to create a tse cngw_certificate
---

# tencentcloud_tse_cngw_certificate

Provides a resource to create a tse cngw_certificate

## Example Usage

```hcl
resource "tencentcloud_tse_cngw_certificate" "cngw_certificate" {
  gateway_id   = "gateway-ddbb709b"
  bind_domains = ["example1.com"]
  cert_id      = "vYSQkJ3K"
  name         = "xxx1"
}
```

## Argument Reference

The following arguments are supported:

* `bind_domains` - (Required, Set: [`String`]) Domains of the binding.
* `cert_id` - (Required, String, ForceNew) Certificate ID of ssl platform.
* `gateway_id` - (Required, String, ForceNew) Gateway ID.
* `name` - (Optional, String) Certificate name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `crt` - Pem format of certificate.
* `key` - Private key of certificate.


## Import

tse cngw_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_certificate.cngw_certificate gatewayId#Id
```

