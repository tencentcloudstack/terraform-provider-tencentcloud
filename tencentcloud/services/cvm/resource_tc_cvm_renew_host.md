Provides a resource to create a cvm renew_host

Example Usage

```hcl
resource "tencentcloud_cvm_renew_host" "renew_host" {
  host_id = "xxxxxx"
  host_charge_prepaid {
	period = 1
	renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  }
}
```
