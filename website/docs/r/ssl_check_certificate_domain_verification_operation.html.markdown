---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_check_certificate_domain_verification_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_check_certificate_domain_verification_operation"
description: |-
  Provides a resource to create a ssl Check Certificate Domain Verification
---

# tencentcloud_ssl_check_certificate_domain_verification_operation

Provides a resource to create a ssl Check Certificate Domain Verification

~> **NOTE:** You can customize the maximum timeout time by setting parameter `timeouts`, which defaults to 15 minutes.

## Example Usage

### Check certificate domain

```hcl
resource "tencentcloud_ssl_check_certificate_domain_verification_operation" "example" {
  certificate_id = "6BE701Jx"
}
```

### Check certificate domain and set the maximum timeout period

```hcl
resource "tencentcloud_ssl_check_certificate_domain_verification_operation" "example" {
  certificate_id = "6BE701Jx"

  timeouts {
    create = "30m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String, ForceNew) The certificate ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `verification_results` - Domain name verification results.
  * `ca_check` - CA inspection results.
  * `check_value` - Detected values.
  * `domain` - Domain name.
  * `frequently` - Whether frequent requests.
  * `issued` - Whether issued.
  * `local_check_fail_reason` - Check the reason for the failure.
  * `local_check` - Local inspection results.
  * `verify_type` - Domain Verify Type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `15m`) Used when creating the resource.

