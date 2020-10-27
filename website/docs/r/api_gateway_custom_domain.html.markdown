---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_custom_domain"
sidebar_current: "docs-tencentcloud-resource-api_gateway_custom_domain"
description: |-
  Use this resource to create custom domain of API gateway.
---

# tencentcloud_api_gateway_custom_domain

Use this resource to create custom domain of API gateway.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_custom_domain" "foo" {
  service_id         = "service-ohxqslqe"
  sub_domain         = "tic-test.dnsv1.com"
  protocol           = "http"
  net_type           = "OUTER"
  is_default_mapping = "false"
  default_domain     = "service-ohxqslqe-1259649581.gz.apigw.tencentcs.com"
  path_mappings      = ["/good#test", "/root#release"]
}
```

## Argument Reference

The following arguments are supported:

* `default_domain` - (Required) Default domain name.
* `net_type` - (Required) Network type. Valid values: `OUTER`, `INNER`.
* `protocol` - (Required) Protocol supported by service. Valid values: `http`, `https`, `http&https`.
* `service_id` - (Required, ForceNew) Unique service ID.
* `sub_domain` - (Required) Custom domain name to be bound.
* `certificate_id` - (Optional) Unique certificate ID of the custom domain name to be bound. The certificate can be uploaded if Protocol is `https` or `http&https`.
* `is_default_mapping` - (Optional) Whether the default path mapping is used. The default value is true. If the value is false, the custom path mapping will be used and PathMappingSet will be required in this case.
* `path_mappings` - (Optional) Custom domain name path mapping. Valid values: `test`, `prepub`, `release`. Respectively.eg: path#environment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Domain name resolution status. True: success; False: failure.


