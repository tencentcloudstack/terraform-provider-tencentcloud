---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_cloud_resource_attachment"
sidebar_current: "docs-tencentcloud-resource-kms_cloud_resource_attachment"
description: |-
  Provides a resource to create a kms cloud_resource_attachment
---

# tencentcloud_kms_cloud_resource_attachment

Provides a resource to create a kms cloud_resource_attachment

## Example Usage

```hcl
resource "tencentcloud_kms_cloud_resource_attachment" "example" {
  key_id      = "72688f39-1fe8-11ee-9f1a-525400cf25a4"
  product_id  = "mysql"
  resource_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `key_id` - (Required, String, ForceNew) CMK unique identifier.
* `product_id` - (Required, String, ForceNew) A unique identifier for the cloud product.
* `resource_id` - (Required, String, ForceNew) The resource/instance ID of the cloud product.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `alias` - Alias.
* `description` - Description.
* `key_state` - Key state.
* `key_usage` - Key usage.
* `owner` - owner.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [BindCloudResource](https://cloud.tencent.com/document/api/573/46918) - This API is used to bind a key with a Tencent Cloud resource. If the key has been set to be expired automatically, the setting will be canceled to ensure that the key will not be invalid automatically. If the key and the resource has already been bound, the call will still be successful.
* Read: [DescribeKey](https://cloud.tencent.com/document/api/573/34428) - This API is used to get the attribute details of the CMK with a specified `KeyId`.
* Delete: [UnbindCloudResource](https://cloud.tencent.com/document/api/573/46917) - This API is used to unbind a key with a Tencent Cloud resource, indicating that the Tencent Cloud resource will not use the key any longer.

## Import

kms cloud_resource_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_kms_cloud_resource_attachment.example 72688f39-1fe8-11ee-9f1a-525400cf25a4#mysql#cdb-fitq5t9h
```

