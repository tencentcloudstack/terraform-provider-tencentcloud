Provides a mysql instance resource to create CDB dr(disaster recovery) instance.

~> **NOTE:** Field `charge_type` only supports modification from `POSTPAID` to `PREPAID`. And the default renewal period is 1 month. and you can also use the `prepaid_period` field to customize the renewal period.

Example Usage

Create POSTPAID dr instance

```hcl
resource "tencentcloud_mysql_dr_instance" "example" {
  master_instance_id = "cdb-3kwa3gfj"
  master_region      = "ap-guangzhou"
  auto_renew_flag    = 0
  availability_zone  = "ap-guangzhou-6"
  charge_type        = "POSTPAID"
  cpu                = 4
  device_type        = "UNIVERSAL"
  first_slave_zone   = "ap-guangzhou-7"
  instance_name      = "tf-example"
  mem_size           = 8000
  project_id         = 0
  security_groups = [
    "sg-e6a8xxib",
  ]
  slave_deploy_mode = 1
  slave_sync_mode   = 0
  subnet_id         = "subnet-hhi88a58"
  volume_size       = 100
  vpc_id            = "vpc-i5yyodl9"
  intranet_port     = 3360
  tags = {
    createBy = "Terraform"
  }
}
```

Create PREPAID dr instance

```hcl
resource "tencentcloud_mysql_dr_instance" "example" {
  master_instance_id = "cdb-3kwa3gfj"
  master_region      = "ap-guangzhou"
  availability_zone  = "ap-guangzhou-6"
  charge_type        = "PREPAID"
  prepaid_period     = 1
  auto_renew_flag    = 1
  cpu                = 4
  device_type        = "UNIVERSAL"
  first_slave_zone   = "ap-guangzhou-7"
  instance_name      = "tf-example"
  mem_size           = 8000
  project_id         = 0
  security_groups = [
    "sg-e6a8xxib",
  ]
  slave_deploy_mode = 1
  slave_sync_mode   = 0
  subnet_id         = "subnet-hhi88a58"
  volume_size       = 100
  vpc_id            = "vpc-i5yyodl9"
  intranet_port     = 3360
  tags = {
    createBy = "Terraform"
  }
}
```

Import

CDB dr(disaster recovery) instancecan be imported using the id, e.g.

```
terraform import tencentcloud_mysql_dr_instance.example cdb-bcet7sdb
```
