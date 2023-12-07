Use this resource to create postgresql readonly instance.

Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_instance" "foo" {
  auto_renew_flag       = 0
  db_version            = "10.4"
  instance_charge_type  = "POSTPAID_BY_HOUR"
  master_db_instance_id = "postgres-j4pm65id"
  memory                = 4
  name                  = "hello"
  need_support_ipv6     = 0
  project_id            = 0
  security_groups_ids   = [
    "sg-fefj5n6r",
  ]
  storage               = 250
  subnet_id             = "subnet-enm92y0m"
  vpc_id                = "vpc-86v957zb"
  read_only_group_id    = tencentcloud_postgresql_readonly_group.new_ro_group.id
}

  resource "tencentcloud_postgresql_readonly_group" "new_ro_group" {
	master_db_instance_id = local.pgsql_id
	name = "tf_ro_group_test_new"
	project_id = 0
	vpc_id  = local.vpc_id
	subnet_id 	= local.subnet_id
	replay_lag_eliminate = 1
	replay_latency_eliminate =  1
	max_replay_lag = 100
	max_replay_latency = 512
	min_delay_eliminate_reserve = 1
  }
```

Import

postgresql readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_postgresql_readonly_instance.foo instance_id
```