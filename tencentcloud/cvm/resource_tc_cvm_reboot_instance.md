Provides a resource to create a cvm reboot_instance

Example Usage

```hcl
resource "tencentcloud_cvm_reboot_instance" "reboot_instance" {
  instance_id = "ins-f9jr4bd2"
  stop_type = "SOFT_FIRST"
}
```