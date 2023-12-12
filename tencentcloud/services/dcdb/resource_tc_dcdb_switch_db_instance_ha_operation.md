Provides a resource to create a dcdb switch_db_instance_ha_operation

Example Usage

```hcl
resource "tencentcloud_dcdb_switch_db_instance_ha_operation" "switch_operation" {
  instance_id = local.dcdb_id
  zone = "ap-guangzhou-4" //3 to 4
}
```