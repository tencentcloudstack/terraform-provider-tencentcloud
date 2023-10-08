---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_deploy_certificate_record_rollback_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_deploy_certificate_record_rollback_operation"
description: |-
  Provides a resource to create a ssl deploy_certificate_record_rollback
---

# tencentcloud_ssl_deploy_certificate_record_rollback_operation

Provides a resource to create a ssl deploy_certificate_record_rollback

## Example Usage

```hcl
resource "tencentcloud_ssl_deploy_certificate_record_rollback_operation" "deploy_certificate_record_rollback" {
  deploy_record_id = 35471
}
```

## Argument Reference

The following arguments are supported:

* `deploy_record_id` - (Optional, Int, ForceNew) Deployment record ID to be rollback.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ssl deploy_certificate_record_rollback can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_deploy_certificate_record_rollback_operation.deploy_certificate_record_rollback deploy_certificate_record_rollback_id
```

