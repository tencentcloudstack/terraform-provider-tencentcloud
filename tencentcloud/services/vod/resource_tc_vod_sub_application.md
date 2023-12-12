Provide a resource to create a VOD sub application.

Example Usage

```hcl
resource  "tencentcloud_vod_sub_application" "foo" {
  name = "foo"
  status = "On"
  description = "this is sub application"
}
```

Import

VOD super player config can be imported using the name+, e.g.

```
$ terraform import tencentcloud_vod_sub_application.foo name+"#"+id
```