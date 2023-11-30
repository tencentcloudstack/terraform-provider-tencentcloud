Provides a resource to create a cls config attachment

Example Usage

```hcl
resource "tencentcloud_cls_config_attachment" "attach" {
  config_id = tencentcloud_cls_config.config.id
  group_id = "27752a9b-9918-440a-8ee7-9c84a14a47ed"
}

Import

cls config_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cls_config_attachment.attach config_id#group_id
```

```