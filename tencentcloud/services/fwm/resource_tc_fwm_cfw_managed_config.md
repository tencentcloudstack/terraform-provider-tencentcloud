Provides a resource to manage the cloud firewall managed switch of an organization member for the fwm (Cloud Firewall Managed) product.

Example Usage

```hcl
resource "tencentcloud_fwm_cfw_managed_config" "example" {
  member_id   = "mem-tencent-84035535ae473b76"
  cfw_managed = 1
}
```

Import

fwm cfw managed can be imported using the composite id, e.g.

```
terraform import tencentcloud_fwm_cfw_managed_config.example mem-tencent-84035535ae473b76
```
