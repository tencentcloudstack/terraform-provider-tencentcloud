Provides a resource to create a dcdb flush_binlog_operation

Example Usage

```hcl
resource "tencentcloud_dcdb_flush_binlog_operation" "flush_operation" {
  instance_id = local.dcdb_id
}
```