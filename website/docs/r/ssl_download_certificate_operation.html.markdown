---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_download_certificate_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_download_certificate_operation"
description: |-
  Provides a resource to create a ssl download_certificate
---

# tencentcloud_ssl_download_certificate_operation

Provides a resource to create a ssl download_certificate

## Example Usage

```hcl
resource "tencentcloud_ssl_download_certificate_operation" "download_certificate" {
  certificate_id = "8x1eUSSl"
  output_path    = "./"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String, ForceNew) Certificate ID.
* `output_path` - (Required, String, ForceNew) Certificate ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ssl download_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_download_certificate_operation.download_certificate download_certificate_id
```

