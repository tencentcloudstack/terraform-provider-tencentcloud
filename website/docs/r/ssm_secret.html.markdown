---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_secret"
sidebar_current: "docs-tencentcloud-resource-ssm_secret"
description: |-
  Provide a resource to create a SSM secret.
---

# tencentcloud_ssm_secret

Provide a resource to create a SSM secret.

## Example Usage

### Create user defined secret

```hcl
resource "tencentcloud_ssm_secret" "example" {
  secret_name             = "tf-example"
  description             = "desc."
  is_enabled              = true
  recovery_window_in_days = 0

  tags = {
    createBy = "terraform"
  }
}
```

### Create redis secret

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 8
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[3].zone
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
}

resource "tencentcloud_redis_instance" "example" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[3].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[3].type_id
  password           = "Qwer@234"
  mem_size           = data.tencentcloud_redis_zone_config.zone.list[3].mem_sizes[0]
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[3].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[3].redis_replicas_nums[0]
  name               = "tf_example"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_ssm_secret" "example" {
  secret_name = "tf-example"
  description = "redis desc."
  is_enabled  = true
  secret_type = 4
  additional_config = jsonencode(
    {
      "Region" : "ap-guangzhou"
      "Privilege" : "r",
      "InstanceId" : tencentcloud_redis_instance.example.id
      "ReadonlyPolicy" : ["master"],
      "Remark" : "for tf test"
    }
  )
  tags = {
    createdBy = "terraform"
  }
  recovery_window_in_days = 0
}
```

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, String, ForceNew) Name of secret which cannot be repeated in the same region. The maximum length is 128 bytes. The name can only contain English letters, numbers, underscore and hyphen '-'. The first character must be a letter or number.
* `additional_config` - (Optional, String) Additional config for specific secret types in JSON string format.
* `description` - (Optional, String) Description of secret. The maximum is 2048 bytes.
* `is_enabled` - (Optional, Bool) Specify whether to enable secret. Default value is `true`.
* `kms_key_id` - (Optional, String, ForceNew) KMS keyId used to encrypt secret. If it is empty, it means that the CMK created by SSM for you by default is used for encryption. You can also specify the KMS CMK created by yourself in the same region for encryption.
* `recovery_window_in_days` - (Optional, Int) Specify the scheduled deletion date. Default value is `0` that means to delete immediately. 1-30 means the number of days reserved, completely deleted after this date.
* `secret_type` - (Optional, Int) Type of secret. `0`: user-defined secret. `4`: redis secret. Default is `0`.
* `tags` - (Optional, Map) Tags of secret.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Status of secret.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSecret](https://cloud.tencent.com/document/api/1140/40529) - This API is used to create a KMS-encrypted Secret. You can create and store up to 1,000 Secrets in each region.
* Create: [DeleteSecretVersion](https://cloud.tencent.com/document/api/1140/40527) - This API is used to directly delete a single credential version under the specified credential. The deletion takes effect immediately, and the credential version in all status can be deleted.
* Create: [DescribeSecret](https://cloud.tencent.com/document/api/1140/40526) - This API is used to obtain the detailed attribute information of a Secret.
* Create: [DisableSecret](https://cloud.tencent.com/document/api/1140/40525) - This API is used to disable a Secret and will change its status to `Disabled`. The plaintext of a disabled Secret cannot be obtained through APIs.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Read: [DescribeSecret](https://cloud.tencent.com/document/api/1140/40526) - This API is used to obtain the detailed attribute information of a Secret.
* Update: [DescribeSecret](https://cloud.tencent.com/document/api/1140/40526) - This API is used to obtain the detailed attribute information of a Secret.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [UpdateDescription](https://cloud.tencent.com/document/api/1140/40516) - This API is used to update the description of a Secret. This API can only update Secrets in `Enabled` or `Disabled` status.
* Delete: [DeleteSecret](https://cloud.tencent.com/document/api/1140/40528) - This API is used to delete a Secret. You can set whether to delete the Secret immediately or on schedule using the `RecoveryWindowInDays` parameter. For a Secret to be deleted on schedule, its status will be `PendingDelete` before the scheduled deletion time. You can use `RestoreSecret` to restore a deleted Secret during this time. A deleted Secret will not be restorable after the scheduled deletion time. A Secret can only be deleted after being disabled using `DisableSecret`.
* Delete: [DescribeSecret](https://cloud.tencent.com/document/api/1140/40526) - This API is used to obtain the detailed attribute information of a Secret.
* Delete: [DisableSecret](https://cloud.tencent.com/document/api/1140/40525) - This API is used to disable a Secret and will change its status to `Disabled`. The plaintext of a disabled Secret cannot be obtained through APIs.

## Import

SSM secret can be imported using the secretName, e.g.
```
$ terraform import tencentcloud_ssm_secret.foo test
```

