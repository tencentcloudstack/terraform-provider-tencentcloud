Provides a resource to create a cvm security_group_attachment

Example Usage

```hcl
resource "tencentcloud_cvm_security_group_attachment" "security_group_attachment" {
  security_group_id = "sg-xxxxxxx"
  instance_id = "ins-xxxxxxxx"
}
```

Import

cvm security_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_security_group_attachment.security_group_attachment ${instance_id}#${security_group_id}
```