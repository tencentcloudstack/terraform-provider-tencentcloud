Provides a resource to manage the CAM password rules configuration (UpdatePasswordRules / GetPasswordRules).

This is a singleton, account-wide configuration resource built on the terraform-plugin-framework implementation. Create and update both call the `UpdatePasswordRules` API; read calls `GetPasswordRules`. Destroying this resource only removes it from state, the cloud-side password rules are left unchanged.

-> **Note:** The CAM password rules are a single account-wide configuration. If the rules look like they have already been changed outside of Terraform, the create emits a warning and applying this configuration will overwrite them; run `terraform import tencentcloud_cam_update_password_rules.example cam-password-rules` and inspect the plan first if the existing values should be kept. Every attribute inside `rules` is optional; a value that is left unset keeps its current cloud-side setting.

-> **Note:** `must_contain` is a literal character set, not an enum: write the characters the password must contain, e.g. `"Aa1!"` (upper, lower, digit, special). The console re-orders and normalizes this value (e.g. `"Aa1!"` is stored and read back as `"1!aA"`), and invalid literals are rejected by the API with `InvalidParameter.PasswordRuleError`.

Example Usage

```hcl
resource "tencentcloud_cam_update_password_rules" "example" {
  rules {
    minimum_length                   = 8
    must_contain                     = "Aa1!"
    force_password_change            = 90
    reuse_password_limit             = 5
    retry_password_limit             = 10
    password_expiration_invalidation = 1
  }
}
```

Import

CAM password rules can be imported using the constant id `cam-password-rules`, e.g.

```
terraform import tencentcloud_cam_update_password_rules.example cam-password-rules
```
