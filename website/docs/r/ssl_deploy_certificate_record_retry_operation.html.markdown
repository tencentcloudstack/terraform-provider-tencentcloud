---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_deploy_certificate_record_retry_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_deploy_certificate_record_retry_operation"
description: |-
  Provides a resource to create a ssl deploy_certificate_record_retry
---

# tencentcloud_ssl_deploy_certificate_record_retry_operation

Provides a resource to create a ssl deploy_certificate_record_retry

## Example Usage

```hcl
resource "tencentcloud_ssl_deploy_certificate_record_retry_operation" "deploy_certificate_record_retry" {
  deploy_record_id = 35474
}
```

## Argument Reference

The following arguments are supported:

* `deploy_record_detail_id` - (Optional, Int, ForceNew) Deployment record details ID to be retried.
* `deploy_record_id` - (Optional, Int, ForceNew) Deployment record ID to be retried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ssl deploy_certificate_record_retry can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_deploy_certificate_record_retry_operation.deploy_certificate_record_retry deploy_certificate_record_retry_id
```

