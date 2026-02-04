---
subcategory: "DNSPod"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_subdomain_validate_txt_value_operation"
sidebar_current: "docs-tencentcloud-resource-subdomain_validate_txt_value_operation"
description: |-
  Provides a resource to create a dnspod subdomain_validate_txt_value_operation
---

# tencentcloud_subdomain_validate_txt_value_operation

Provides a resource to create a dnspod subdomain_validate_txt_value_operation

## Example Usage

```hcl
resource "tencentcloud_subdomain_validate_txt_value_operation" "subdomain_validate_txt_value_operation" {
  domain_zone = "www.iac-tf.cloud"
}
```

## Argument Reference

The following arguments are supported:

* `domain_zone` - (Required, String, ForceNew) The subdomain to add Zone domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `domain` - The domain name for which TXT records need to be added.
* `record_type` - Record types need to be added.
* `subdomain` - Host records that need to be added to TXT records.
* `value` - The record value of the TXT record needs to be added.


