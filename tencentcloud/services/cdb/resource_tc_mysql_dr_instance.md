Provides a mysql instance resource to create read-only database instances.

Example Usage

```hcl
resource "tencentcloud_mysql_dr_instance" "mysql_dr" {
  master_instance_id = "cdb-adjdu3t5"
  master_region      = "ap-guangzhou"
  auto_renew_flag    = 0
  availability_zone  = "ap-shanghai-3"
  charge_type        = "POSTPAID"
  cpu                = 4
  device_type        = "UNIVERSAL"
  first_slave_zone   = "ap-shanghai-4"
  instance_name      = "mysql-dr-test-up"
  mem_size           = 8000
  prepaid_period     = 1
  project_id         = 0
  security_groups = [
    "sg-q4d821qk",
  ]
  slave_deploy_mode = 1
  slave_sync_mode   = 0
  subnet_id         = "subnet-5vfntba5"
  volume_size       = 100
  vpc_id            = "vpc-h6s1s3aa"
  intranet_port     = 3360
  tags = {
    test = "test-tf"
  }
}
```
Import

mysql dr database instances can be imported using the id, e.g.
```
terraform import tencentcloud_mysql_dr_instance.mysql_dr cdb-bcet7sdb
```