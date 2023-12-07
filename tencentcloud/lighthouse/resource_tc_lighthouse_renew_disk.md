Provides a resource to create a lighthouse renew_disk

Example Usage

```hcl
resource "tencentcloud_lighthouse_renew_disk" "renew_disk" {
  disk_id = "lhdisk-xxxxxx"
  renew_disk_charge_prepaid {
	period = 1
	renew_flag = "NOTIFY_AND_AUTO_RENEW"
	time_unit = "m"
  }
  auto_voucher = true
}
```