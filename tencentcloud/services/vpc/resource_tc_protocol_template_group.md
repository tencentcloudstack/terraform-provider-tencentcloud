Provides a resource to manage protocol template group.

Example Usage

```hcl
resource "tencentcloud_protocol_template_group" "foo" {
  name                = "group-test"
  template_ids = ["ipl-axaf24151","ipl-axaf24152"]
}
```

Import

Protocol template group can be imported using the id, e.g.

```
$ terraform import tencentcloud_protocol_template_group.foo ppmg-0np3u974
```