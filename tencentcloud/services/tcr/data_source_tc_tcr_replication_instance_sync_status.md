Use this data source to query detailed information of tcr replication_instance_sync_status

Example Usage

```hcl
data "tencentcloud_tcr_replication_instance_sync_status" "sync_status" {
  registry_id             = local.src_registry_id
  replication_registry_id = local.dst_registry_id
  replication_region_id   = local.dst_region_id
  show_replication_log    = false
}
```