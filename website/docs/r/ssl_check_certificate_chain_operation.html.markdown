---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_check_certificate_chain_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_check_certificate_chain_operation"
description: |-
  Provides a resource to create a ssl check_certificate_chain
---

# tencentcloud_ssl_check_certificate_chain_operation

Provides a resource to create a ssl check_certificate_chain

## Example Usage

```hcl
resource "tencentcloud_ssl_check_certificate_chain_operation" "check_certificate_chain" {
  certificate_chain = "-----BEGIN CERTIFICATE--·····---END CERTIFICATE-----"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_chain` - (Required, String, ForceNew) The certificate chain to check.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ssl check_certificate_chain can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_check_certificate_chain_operation.check_certificate_chain check_certificate_chain_id
```

