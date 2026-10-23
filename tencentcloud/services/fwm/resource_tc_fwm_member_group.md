Provides a resource to manage which FWM (Cloud Firewall Management) organization account group a member belongs to.

~> NOTE: When `group_id` is non-empty, the member must have been managed by cloud firewall (a `tencentcloud_fwm_cfw_managed` resource with `cfw_managed = 1`), so the example below uses `depends_on` to ensure the member is managed first. There is a reverse coupling with `tencentcloud_fwm_cfw_managed`: disabling management (cfw_managed 1->0) will automatically move the member out of the group on the server side. If you manage both resources at the same time, set `group_id` to empty before disabling management, otherwise unexpected state drift may occur.

Example Usage

Move a managed member into an account group

```hcl
resource "tencentcloud_fwm_cfw_managed" "example" {
  member_id   = "mem-xxxx"
  cfw_managed = 1
}

resource "tencentcloud_fwm_account_group" "example" {
  group_name = "tf-example-group"
}

resource "tencentcloud_fwm_member_group" "example" {
  member_id = "mem-xxxx"
  app_id    = "1234567890"
  group_id  = tencentcloud_fwm_account_group.example.id

  depends_on = [tencentcloud_fwm_cfw_managed.example]
}
```

Move a member out of an account group

```hcl
resource "tencentcloud_fwm_member_group" "example" {
  member_id = "mem-xxxx"
  app_id    = "1234567890"
  group_id  = ""
}
```

Import

FWM member group resource can be imported using the two-segment composite id, e.g. `app_id#member_id`, where the first segment is the member AppId used to locate the member and the second segment is the member MemberId.

```sh
terraform import tencentcloud_fwm_member_group.example 1234567890#mem-xxxx
```