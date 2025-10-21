---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_deploy_certificate_instance_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_deploy_certificate_instance_operation"
description: |-
  Provides a resource to create a ssl deploy_certificate_instance
---

# tencentcloud_ssl_deploy_certificate_instance_operation

Provides a resource to create a ssl deploy_certificate_instance

## Example Usage

```hcl
resource "tencentcloud_ssl_deploy_certificate_instance_operation" "deploy_certificate_instance" {
  certificate_id   = "8x1eUSSl"
  instance_id_list = ["cdndomain1.example.com|on", "cdndomain1.example.com|off"]
  resource_type    = "cdn"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String, ForceNew) ID of the certificate to be deployed.
* `instance_id_list` - (Required, Set: [`String`], ForceNew) Need to deploy instance list.
* `resource_type` - (Optional, String, ForceNew) Deployed cloud resource type.
* `status` - (Optional, Int, ForceNew) Deployment cloud resource status: Live: -1: The domain name is not associated with a certificate.1:  Domain name https is enabled.0:  Domain name https is closed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ssl deploy_certificate_instance can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_deploy_certificate_instance_operation.deploy_certificate_instance deploy_certificate_instance_id
```

