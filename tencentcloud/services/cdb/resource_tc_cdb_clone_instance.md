Provides a resource to create a CDB clone instance from an existing source instance, optionally rolling back to a specified time point or backup set.

Example Usage

Clone instance with rollback time

```hcl
resource "tencentcloud_cdb_clone_instance" "example" {
  instance_id            = "cdb-c1nl9rpv"
  specified_rollback_time = "2024-01-01 12:00:00"
  memory                 = 4000
  volume                 = 200
  cpu                    = 2
  instance_name          = "tf-clone-example"
  uniq_vpc_id            = "vpc-i5yyodl9"
  uniq_subnet_id         = "subnet-hhi88a58"
  protect_mode           = 0
  deploy_mode            = 0
  slave_zone             = "ap-guangzhou-3"
  device_type            = "UNIVERSAL"
  project_id             = 0
  pay_type               = "USED_PAID"
  zone                   = "ap-guangzhou-3"
}
```

Clone instance with backup id

```hcl
resource "tencentcloud_cdb_clone_instance" "example" {
  instance_id        = "cdb-c1nl9rpv"
  specified_backup_id = 1000001
  memory              = 4000
  volume              = 200
  cpu                 = 2
  instance_name       = "tf-clone-example"
  uniq_vpc_id         = "vpc-i5yyodl9"
  uniq_subnet_id      = "subnet-hhi88a58"
  protect_mode        = 0
  deploy_mode         = 0
  slave_zone          = "ap-guangzhou-3"
  device_type         = "UNIVERSAL"
  project_id          = 0
  pay_type            = "USED_PAID"
  zone                = "ap-guangzhou-3"
}
```

Clone instance with resource tags

```hcl
resource "tencentcloud_cdb_clone_instance" "example" {
  instance_id            = "cdb-c1nl9rpv"
  specified_rollback_time = "2024-01-01 12:00:00"
  memory                 = 4000
  volume                 = 200
  cpu                    = 2
  instance_name          = "tf-clone-example"
  uniq_vpc_id            = "vpc-i5yyodl9"
  uniq_subnet_id         = "subnet-hhi88a58"
  protect_mode           = 0
  deploy_mode            = 0
  slave_zone             = "ap-guangzhou-3"
  device_type            = "UNIVERSAL"
  project_id             = 0
  pay_type               = "USED_PAID"
  zone                   = "ap-guangzhou-3"

  resource_tags {
    key   = "createBy"
    value = "Terraform"
  }
}
```

Import

CDB clone instance can be imported using the cloned instance id, e.g.

```
terraform import tencentcloud_cdb_clone_instance.example cdb-bcet7sdb
```
