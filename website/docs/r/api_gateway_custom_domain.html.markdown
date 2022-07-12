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

* `default_domain` - (Required, String) Default domain name.
* `net_type` - (Required, String) Network type. Valid values: `OUTER`, `INNER`.
* `protocol` - (Required, String) Protocol supported by service. Valid values: `http`, `https`, `http&https`.
* `service_id` - (Required, String, ForceNew) Unique service ID.
* `sub_domain` - (Required, String) Custom domain name to be bound.
* `certificate_id` - (Optional, String) Unique certificate ID of the custom domain name to be bound. You can choose to upload for the `protocol` attribute value `https` or `http&https`.
* `is_default_mapping` - (Optional, Bool) Whether the default path mapping is used. The default value is `true`. When it is `false`, it means custom path mapping. In this case, the `path_mappings` attribute is required.
* `path_mappings` - (Optional, Set: [`String`]) Custom domain name path mapping. The data format is: `path#environment`. Optional values for the environment are `test`, `prepub`, and `release`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Domain name resolution status. `1` means normal analysis, `0` means parsing failed.


