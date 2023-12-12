Provides a reserved instance resource.

~> **NOTE:** Reserved instance cannot be deleted and updated. The reserved instance still exist which can be extracted by reserved_instances data source when reserved instance is destroied.

Example Usage

```hcl
resource "tencentcloud_reserved_instance" "ri" {
  config_id      = "469043dd-28b9-4d89-b557-74f6a8326259"
  instance_count = 2
}
```

Import

Reserved instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_reserved_instance.foo 6cc16e7c-47d7-4fae-9b44-ce5c0f59a920
```