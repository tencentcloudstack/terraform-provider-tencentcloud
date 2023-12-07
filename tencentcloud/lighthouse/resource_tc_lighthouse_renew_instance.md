Provides a resource to create a lighthouse renew_instance

Example Usage

```hcl
resource "tencentcloud_lighthouse_renew_instance" "renew_instance" {
  instance_id =
  instance_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_MANUAL_RENEW"

  }
  renew_data_disk = true
  auto_voucher = false
}
```