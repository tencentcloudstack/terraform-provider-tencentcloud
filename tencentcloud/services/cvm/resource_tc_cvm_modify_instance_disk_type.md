Provides a resource to create a cvm modify_instance_disk_type

Example Usage

```hcl
resource "tencentcloud_cvm_modify_instance_disk_type" "modify_instance_disk_type" {
  instance_id = "ins-r8hr2upy"
  data_disks {
		disk_size = 50
		disk_type = "CLOUD_BASIC"
		disk_id = "disk-hrsd0u81"
		delete_with_instance = true
		snapshot_id = "snap-r9unnd89"
		encrypt = false
		kms_key_id = "kms-abcd1234"
		throughput_performance = 2
		cdc_id = "cdc-b9pbd3px"

  }
  system_disk {
		disk_type = "CLOUD_PREMIUM"
		disk_id = "disk-1drr53sd"
		disk_size = 50
		cdc_id = "cdc-b9pbd3px"

  }
}
```
