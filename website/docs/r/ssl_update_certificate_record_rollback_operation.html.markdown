---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_update_certificate_record_rollback_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_update_certificate_record_rollback_operation"
description: |-
  Provides a resource to create a ssl update_certificate_record_rollback
---

# tencentcloud_ssl_update_certificate_record_rollback_operation

Provides a resource to create a ssl update_certificate_record_rollback

## Example Usage

```hcl
resource "tencentcloud_ssl_update_certificate_record_rollback_operation" "update_certificate_record_rollback" {
  deploy_record_id = "1603"
}
```

## Argument Reference

The following arguments are supported:

* `deploy_record_id` - (Optional, String, ForceNew) Deployment record ID to be rolled back.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ssl update_certificate_record_rollback can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_update_certificate_record_rollback_operation.update_certificate_record_rollback update_certificate_record_rollback_id
```

