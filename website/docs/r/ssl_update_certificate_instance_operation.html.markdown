---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_update_certificate_instance_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_update_certificate_instance_operation"
description: |-
  Provides a resource to create a ssl update_certificate_instance
---

# tencentcloud_ssl_update_certificate_instance_operation

Provides a resource to create a ssl update_certificate_instance

## Example Usage

```hcl
resource "tencentcloud_ssl_update_certificate_instance_operation" "update_certificate_instance" {
  certificate_id     = "8x1eUSSl"
  old_certificate_id = "8xNdi2ig"
  resource_types     = ["cdn"]
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String, ForceNew) Update new certificate ID.
* `old_certificate_id` - (Required, String, ForceNew) Update the original certificate ID.
* `resource_types` - (Required, Set: [`String`], ForceNew) The resource type that needs to be deployed. The parameter value is optional: clb,cdn,waf,live,ddos,teo,apigateway,vod,tke,tcb.
* `resource_types_regions` - (Optional, List, ForceNew) List of regions where cloud resources need to be deployed.

The `resource_types_regions` object supports the following:

* `regions` - (Optional, Set) Region list.
* `resource_type` - (Optional, String) Cloud resource type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ssl update_certificate_instance can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_update_certificate_instance_operation.update_certificate_instance update_certificate_instance_id
```

