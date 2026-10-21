---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_upload_certificate_instance"
sidebar_current: "docs-tencentcloud-resource-ssl_upload_certificate_instance"
description: |-
  Provides a resource to create a SSL upload certificate instance. The `tencentcloud_ssl_upload_certificate_instance` resource uploads new certificate public/private key content and updates associated cloud resources (such as clb) without changing the certificate id.
---

# tencentcloud_ssl_upload_certificate_instance

Provides a resource to create a SSL upload certificate instance. The `tencentcloud_ssl_upload_certificate_instance` resource uploads new certificate public/private key content and updates associated cloud resources (such as clb) without changing the certificate id.

## Example Usage

```hcl
resource "tencentcloud_ssl_upload_certificate_instance" "example" {
  old_certificate_id      = "f8xnynbd"
  resource_types          = ["clb"]
  certificate_public_key  = "-----BEGIN CERTIFICATE-----\nxxx\n-----END CERTIFICATE-----"
  certificate_private_key = "-----BEGIN RSA PRIVATE KEY-----\nxxx\n-----END RSA PRIVATE KEY-----"
  resource_types_regions {
    resource_type = "clb"
    regions       = ["ap-guangzhou"]
  }
}
```

### Upload with certificate id

```hcl
resource "tencentcloud_ssl_upload_certificate_instance" "example" {
  old_certificate_id = "f8xnynbd"
  certificate_id     = "f8xnynbe"
  resource_types     = ["clb"]
}
```

## Argument Reference

The following arguments are supported:

* `certificate_private_key` - (Required, String) Certificate private key. If you upload the certificate public key, CertificateId does not need to be passed.
* `certificate_public_key` - (Required, String) Certificate public key. If you upload the certificate public key, CertificateId does not need to be passed.
* `old_certificate_id` - (Required, String) Update the original certificate ID. The cloud resources bound to this certificate ID are queried, and then these cloud resources are updated with the new certificate.
* `resource_types` - (Required, Set: [`String`]) The resource type that needs to be deployed. The parameter value is optional: clb.
* `certificate_id` - (Optional, String) Certificate ID.
* `resource_types_regions` - (Optional, List) List of regions where cloud resources need to be deployed.

The `resource_types_regions` object supports the following:

* `regions` - (Optional, Set) Region list.
* `resource_type` - (Optional, String) Cloud resource type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `deploy_record_id` - Cloud resource update task ID. 0 means the task is in progress. When DeployRecordId is greater than 0, the task is created successfully.
* `deploy_record_list` - Certificate deployment record list.
  * `create_time` - Deployment time.
  * `id` - Record ID.
  * `old_cert_id` - Original certificate ID.
  * `resource_types` - List of deployed resource types.
  * `status` - Deployment status.
  * `update_time` - Last update time.
* `deploy_status` - Update task creation status. 1 means created successfully; 0 means there is an update task in progress, and no new update task is created.
* `update_sync_progress` - Update asynchronous creation task progress details.
  * `resource_type` - Resource type.
  * `status` - Asynchronous update progress status: 0 pending, 1 processed, 3 processing.
  * `update_sync_progress_regions` - Region result list.
    * `offset_count` - Number of completed executions.
    * `region` - Region.
    * `status` - Asynchronous update progress status: 0 pending, 1 processed, 3 processing.
    * `total_count` - Total count.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [DescribeHostUploadUpdateRecord](https://www.tencentcloud.com/document/api/1007/70505)
  - [DescribeHostUploadUpdateRecordDetail](https://www.tencentcloud.com/document/api/1007/70504)
  - [UploadUpdateCertificateInstance](https://www.tencentcloud.com/document/api/1007/70503)
* Read:
  - [DescribeHostUploadUpdateRecord](https://www.tencentcloud.com/document/api/1007/70505)
  - [DescribeHostUploadUpdateRecordDetail](https://www.tencentcloud.com/document/api/1007/70504)
* Update:
  - [DescribeHostUploadUpdateRecord](https://www.tencentcloud.com/document/api/1007/70505)
  - [DescribeHostUploadUpdateRecordDetail](https://www.tencentcloud.com/document/api/1007/70504)
  - [UploadUpdateCertificateInstance](https://www.tencentcloud.com/document/api/1007/70503)
* Delete: [UploadUpdateCertificateRecordRollback](https://www.tencentcloud.com/document/api/1007/70501)

## Import

SSL upload certificate instance can be imported using the deploy record id, e.g.

```
$ terraform import tencentcloud_ssl_upload_certificate_instance.example 12345
```

