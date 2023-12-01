---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_complete_certificate_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_complete_certificate_operation"
description: |-
  Provides a resource to create a ssl complete_certificate
---

# tencentcloud_ssl_complete_certificate_operation

Provides a resource to create a ssl complete_certificate

## Example Usage

```hcl
resource "tencentcloud_ssl_complete_certificate_operation" "complete_certificate" {
  certificate_id = "9Bfe1IBR"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String, ForceNew) Certificate ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ssl complete_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_complete_certificate_operation.complete_certificate complete_certificate_id
```

