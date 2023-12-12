Provides a resource to manage address template group.

Example Usage

```hcl
resource "tencentcloud_address_template_group" "foo" {
  name                = "group-test"
  template_ids = ["ipl-axaf24151","ipl-axaf24152"]
}
```

Import

Address template group can be imported using the id, e.g.

```
$ terraform import tencentcloud_address_template_group.foo ipmg-0np3u974
```