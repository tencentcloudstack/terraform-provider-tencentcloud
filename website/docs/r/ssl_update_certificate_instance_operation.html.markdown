---
subcategory: "SSL Certificates(ssl)"
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

### Upload certificate

```hcl
resource "tencentcloud_ssl_update_certificate_instance_operation" "update_certificate_instance" {
  old_certificate_id      = "xxx"
  certificate_public_key  = file("xxx.crt")
  certificate_private_key = file("xxx.key")
  repeatable              = true
  resource_types          = ["cdn"]
}
```

## Argument Reference

The following arguments are supported:

* `old_certificate_id` - (Required, String, ForceNew) Update the original certificate ID.
* `resource_types` - (Required, Set: [`String`], ForceNew) The resource type that needs to be deployed. The parameter value is optional: clb, cdn, waf, live, ddos, teo, apigateway, vod, tke, tcb.
* `allow_download` - (Optional, Bool, ForceNew) Whether to allow downloading, if you choose to upload the certificate, you can configure this parameter.
* `certificate_id` - (Optional, String, ForceNew) Update new certificate ID.
* `certificate_private_key` - (Optional, String, ForceNew) Certificate private key. If you upload the certificate public key, CertificateId does not need to be passed.
* `certificate_public_key` - (Optional, String, ForceNew) Certificate public key. If you upload the certificate public key, CertificateId does not need to be passed.
* `expiring_notification_switch` - (Optional, Int, ForceNew) Whether to ignore expiration reminders for old certificates 0: Do not ignore notifications. 1: Ignore the notification and ignore the OldCertificateId expiration reminder.
* `project_id` - (Optional, Int, ForceNew) Project ID, if you choose to upload the certificate, you can configure this parameter.
* `repeatable` - (Optional, Bool, ForceNew) Whether the same certificate is allowed to be uploaded repeatedly. If you choose to upload the certificate, you can configure this parameter.
* `resource_types_regions` - (Optional, List, ForceNew) List of regions where cloud resources need to be deploye.

The `resource_types_regions` object supports the following:

* `regions` - (Optional, Set) Region list.
* `resource_type` - (Optional, String) Cloud resource type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



