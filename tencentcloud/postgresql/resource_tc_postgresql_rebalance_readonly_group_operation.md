Provides a resource to create a postgresql rebalance_readonly_group_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_group" "group_rebalance" {
	master_db_instance_id = local.pgsql_id
	name = "test-pg-readonly-group-rebalance"
	project_id = 0
	vpc_id = "vpc-86v957zb"
	subnet_id = "subnet-enm92y0m"
	replay_lag_eliminate = 1
	replay_latency_eliminate =  1
	max_replay_lag = 100
	max_replay_latency = 512
	min_delay_eliminate_reserve = 1
}

resource "tencentcloud_postgresql_rebalance_readonly_group_operation" "rebalance_readonly_group_operation" {
  read_only_group_id = tencentcloud_postgresql_readonly_group.group_rebalance.id
}
```