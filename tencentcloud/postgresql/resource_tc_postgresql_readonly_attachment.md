Use this resource to create postgresql readonly attachment.

Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_attachment" "attach" {
  db_instance_id = tencentcloud_postgresql_readonly_instance.foo.id
  read_only_group_id = tencentcloud_postgresql_readonly_group.group.id
}
```