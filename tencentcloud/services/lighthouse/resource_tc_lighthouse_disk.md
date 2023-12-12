Provides a resource to create a lighthouse disk

Example Usage

```hcl
resource "tencentcloud_lighthouse_disk" "disk" {
  zone = "ap-hongkong-2"
  disk_size = 20
  disk_type = "CLOUD_SSD"
  disk_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_AUTO_RENEW"
		time_unit = "m"

  }
  disk_name = "test"
}
```