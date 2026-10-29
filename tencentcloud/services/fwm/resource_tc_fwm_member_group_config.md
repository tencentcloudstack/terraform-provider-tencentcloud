Provides a resource to manage which FWM (Cloud Firewall Management) organization account group a member belongs to.

~> NOTE: When `group_id` is non-empty, the member must have been managed by cloud firewall (a `tencentcloud_fwm_cfw_managed` resource with `cfw_managed = 1`), so the example below uses `depends_on` to ensure the member is managed first. There is a reverse coupling with `tencentcloud_fwm_cfw_managed`: disabling management (cfw_managed 1->0) will automatically move the member out of the group on the server side. If you manage both resources at the same time, set `group_id` to empty before disabling management, otherwise unexpected state drift may occur.

Example Usage

```hcl
resource "tencentcloud_fwm_member_group_config" "example" {
  user_member_id = "mem-tencent-26afe1f469a10371"
  group_id       = "acg-rqjwnsnc"
}
```

Import

FWM member group resource can be imported using the id, e.g.

```
terraform import tencentcloud_fwm_member_group_config.example mem-tencent-26afe1f469a10371
```