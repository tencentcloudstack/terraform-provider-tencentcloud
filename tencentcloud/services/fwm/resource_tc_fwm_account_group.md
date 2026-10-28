Provides a resource to create a FWM (Edge firewall account group management) account group

Example Usage

```hcl
resource "tencentcloud_fwm_account_group" "example" {
  group_name  = "tf-example"
  description = "remark."
}
```

Import

FWM account group can be imported using the GroupId, e.g.

```
terraform import tencentcloud_fwm_account_group.example acg-nyalbmbi
```