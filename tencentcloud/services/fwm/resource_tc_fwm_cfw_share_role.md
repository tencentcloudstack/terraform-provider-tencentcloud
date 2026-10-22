Provides a resource to create a fwm cfw_share_role, used to manage the cloud firewall (fwm) share role (sharer/user/none) expectation state of an organization member.

Example Usage

```hcl
resource "tencentcloud_fwm_cfw_share_role" "example" {
  member_id     = "1234567890"
  app_id        = "1234567890"
  cfw_share_role = "none"
}
```

Import

fwm cfw share role can be imported using the appId#memberId, e.g.

```
terraform import tencentcloud_fwm_cfw_share_role.example 1234567890#1234567890
```