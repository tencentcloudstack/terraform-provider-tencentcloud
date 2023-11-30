Use this data source to query detailed information of tcr replication_instance_create_tasks

Example Usage

```hcl
local {
  src_registry_id = local.tcr_id
  dst_registry_id = tencentcloud_tcr_manage_replication_operation.my_replica.destination_registry_id
  dst_region_id   = tencentcloud_tcr_manage_replication_operation.my_replica.destination_region_id
}

data "tencentcloud_tcr_replication_instance_create_tasks" "create_tasks" {
  replication_registry_id = local.dst_registry_id
  replication_region_id   = local.dst_region_id
}
```