Use this resource to create postgresql readonly group.

Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_group" "group" {
  master_db_instance_id = "postgres-gzg9jb2n"
  name = "world"
  project_id = 0
  vpc_id = "vpc-86v957zb"
  subnet_id = "subnet-enm92y0m"
  replay_lag_eliminate = 1
  replay_latency_eliminate =  1
  max_replay_lag = 100
  max_replay_latency = 512
  min_delay_eliminate_reserve = 1
#  security_groups_ids = []
}
```