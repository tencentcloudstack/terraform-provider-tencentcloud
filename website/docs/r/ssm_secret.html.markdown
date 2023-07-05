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
resource "tencentcloud_ssm_secret" "foo" {
  secret_name             = "test"
  description             = "user defined secret"
  recovery_window_in_days = 0
  is_enabled              = true

  tags = {
    test-tag = "test"
  }
}
```

### Create redis secret

```hcl
data "tencentcloud_redis_instances" "instance" {
  zone = "ap-guangzhou-6"
}

resource "tencentcloud_ssm_secret" "secret" {
  secret_name = "for-redis-test"
  description = "redis secret"
  is_enabled  = false

  secret_type = 4
  additional_config = jsonencode(
    {
      "Region" : "ap-guangzhou"
      "Privilege" : "r",
      "InstanceId" : data.tencentcloud_redis_instances.instance.instance_list.0.redis_id
      "ReadonlyPolicy" : ["master"],
      "Remark" : "for tf test"
    }
  )
  tags = {
    test-tag = "test"
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
* `secret_type` - (Optional, Int) Type of secret. `0`: user-defined secret. `4`: redis secret.
* `tags` - (Optional, Map) Tags of secret.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Status of secret.


## Import

SSM secret can be imported using the secretName, e.g.
```
$ terraform import tencentcloud_ssm_secret.foo test
```

