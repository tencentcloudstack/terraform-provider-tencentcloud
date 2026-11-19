Provides a resource to manage the backup encryption config of a Redis instance.

Example Usage

Enable backup encryption with Tencent Cloud managed key

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
  subnet_id         = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_redis_backup_encrypt" "example" {
  instance_id    = tencentcloud_redis_instance.example.id
  backup_encrypt = true
}
```

Enable backup encryption with user-defined KMS key

```hcl
resource "tencentcloud_redis_backup_encrypt" "example" {
  instance_id      = tencentcloud_redis_instance.example.id
  backup_encrypt    = true
  encrypt_key_type = 1
  kms_key_id        = "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx"
  kms_region        = "ap-guangzhou"
}
```

Import

Redis backup encrypt config can be imported, e.g.

```
$ terraform import tencentcloud_redis_backup_encrypt.example crs-cqdfdzvt
```