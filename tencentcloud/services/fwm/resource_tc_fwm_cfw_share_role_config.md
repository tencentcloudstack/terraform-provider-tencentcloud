Provides a resource to create a fwm cfw_share_role, used to manage the cloud firewall (fwm) share role (sharer/user/none) expectation state of an organization member.

Example Usage

Configure the share role as user

```hcl
resource "tencentcloud_fwm_cfw_share_role_config" "example" {
  cfw_share_role   = "user"
  user_member_id   = "mem-tencent-e46e42e79411caf0"
  sharer_member_id = "mem-tencent-d55a312b431535c3"
}
```

Configure the share role as none

```hcl
resource "tencentcloud_fwm_cfw_share_role_config" "example" {
  cfw_share_role = "none"
  user_member_id = "mem-tencent-e46e42e79411caf0"
}
```

Import

fwm cfw share role can be imported using id, e.g.

```
terraform import tencentcloud_fwm_cfw_share_role_config.example mem-tencent-e46e42e79411caf0
```