Provides a resource to manage the cloud firewall managed switch of an organization member for the fwm (Cloud Firewall Managed) product.

Example Usage

```hcl
resource "tencentcloud_fwm_cfw_managed" "example" {
  app_id      = "1234567890"
  member_id   = "100012345678"
  cfw_managed = 1
}
```

Import

fwm cfw managed can be imported using the composite id, e.g. `app_id#member_id`.

```
terraform import tencentcloud_fwm_cfw_managed.example 1234567890#100012345678
```
