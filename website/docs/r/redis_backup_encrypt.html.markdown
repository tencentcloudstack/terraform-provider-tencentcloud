---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_backup_encrypt"
sidebar_current: "docs-tencentcloud-resource-redis_backup_encrypt"
description: |-
  Provides a resource to manage the backup encryption config of a Redis instance.
---

# tencentcloud_redis_backup_encrypt

Provides a resource to manage the backup encryption config of a Redis instance.

## Example Usage

### Enable backup encryption with Tencent Cloud managed key

```hcl
resource "tencentcloud_redis_instance" "example" {
  availability_zone  = "ap-guangzhou-3"
  type_id            = 7
  password           = "Password@123"
  mem_size           = 8192
  redis_shard_num    = 1
  redis_replicas_num = 1
  name               = "tf_example"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_redis_backup_encrypt" "example" {
  instance_id    = tencentcloud_redis_instance.example.id
  backup_encrypt = true
}
```

### Enable backup encryption with user-defined KMS key

```hcl
resource "tencentcloud_redis_backup_encrypt" "example" {
  instance_id      = tencentcloud_redis_instance.example.id
  backup_encrypt   = true
  encrypt_key_type = 1
  kms_key_id       = "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx"
  kms_region       = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `backup_encrypt` - (Required, Bool) Whether to enable backup encryption. true: enable; false: disable.
* `instance_id` - (Required, String, ForceNew) The ID of redis instance.
* `encrypt_key_type` - (Optional, Int) Encryption key type. 0: Tencent Cloud managed (default); 1: user-defined KMS key.
* `kms_key_id` - (Optional, String) The KMS key ID.
* `kms_region` - (Optional, String) The KMS region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [DescribeBackupEncryptConfig](https://www.tencentcloud.com/document/api/239/85025)
  - [ModifyBackupEncrypt](https://www.tencentcloud.com/document/api/239/85024)
* Read: [DescribeBackupEncryptConfig](https://www.tencentcloud.com/document/api/239/85025)
* Update:
  - [DescribeBackupEncryptConfig](https://www.tencentcloud.com/document/api/239/85025)
  - [ModifyBackupEncrypt](https://www.tencentcloud.com/document/api/239/85024)

## Import

Redis backup encrypt config can be imported, e.g.

```
$ terraform import tencentcloud_redis_backup_encrypt.example crs-cqdfdzvt
```

