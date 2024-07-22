---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_transparent_data_encryption"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance_transparent_data_encryption"
description: |-
  Provides a resource to enable mongodb transparent data encryption
---

# tencentcloud_mongodb_instance_transparent_data_encryption

Provides a resource to enable mongodb transparent data encryption

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance_transparent_data_encryption" "encryption" {
  instance_id = "cmgo-xxxxxx"
  kms_region  = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID, for example: cmgo-p8vn ***. Currently supported general versions include: 4.4 and 5.0, but the cloud disk version is not currently supported.
* `kms_region` - (Required, String) The region where the Key Management Service (KMS) serves, such as ap-shanghai.
* `key_id` - (Optional, String) Key ID. If this parameter is not set and the specific key ID is not specified, Tencent Cloud will automatically generate the key and this key will be beyond the control of Terraform.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `key_info_list` - List of bound keys.
  * `create_time` - Instance and key binding time.
  * `key_id` - Master Key ID.
  * `key_name` - Master key name.
  * `key_origin` - Key origin.
  * `key_usage` - Purpose of the key.
  * `status` - Key status.
* `transparent_data_encryption_status` - Represents whether transparent encryption is turned on. Valid values:
- close: Not opened;
- open: It has been opened.


## Import

mongodb transparent data encryption can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance_transparent_data_encryption.encryption ${instanceId}
```

