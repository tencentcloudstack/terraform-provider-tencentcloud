---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_srv_connection"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance_srv_connection"
description: |-
  Provides a resource to manage MongoDB instance SRV connection URL configuration.
---

# tencentcloud_mongodb_instance_srv_connection

Provides a resource to manage MongoDB instance SRV connection URL configuration.

## Example Usage

### Enable SRV connection with default domain

```hcl
resource "tencentcloud_mongodb_instance_srv_connection" "example" {
  instance_id = "cmgo-p8vnipr5"
}

output "domain" {
  value = tencentcloud_mongodb_instance_srv_connection.example.domain
}
```

### Enable SRV connection with custom domain

```hcl
resource "tencentcloud_mongodb_instance_srv_connection" "example" {
  instance_id = "cmgo-p8vnipr5"
  domain      = "exampleDomain"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) MongoDB instance ID, for example: cmgo-p8vnipr5.
* `domain` - (Optional, String) Custom domain for SRV connection. If not specified during creation, the system will use a default domain. After creation, this field will be populated with the actual domain. To set or modify a custom domain, use this field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeAsyncRequestInfo](https://cloud.tencent.com/document/api/240/46601)
* Create: [EnableSRVConnectionUrl](https://cloud.tencent.com/document/api/240/127855)
* Create: [ModifySRVConnectionUrl](https://cloud.tencent.com/document/api/240/127854)
* Read: [DescribeSRVConnectionDomain](https://cloud.tencent.com/document/api/240/127857)
* Update: [DescribeAsyncRequestInfo](https://cloud.tencent.com/document/api/240/46601)
* Update: [ModifySRVConnectionUrl](https://cloud.tencent.com/document/api/240/127854)
* Delete: [DescribeAsyncRequestInfo](https://cloud.tencent.com/document/api/240/46601)
* Delete: [DisableSRVConnectionUrl](https://cloud.tencent.com/document/api/240/127856)

## Import

MongoDB instance SRV connection can be imported using the instance id, e.g.

```
terraform import tencentcloud_mongodb_instance_srv_connection.example cmgo-p8vnipr5
```

