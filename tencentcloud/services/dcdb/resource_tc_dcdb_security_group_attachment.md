Provides a resource to create a dcdb security_group_attachment

Example Usage

```hcl
resource "tencentcloud_dcdb_security_group_attachment" "security_group_attachment" {
  security_group_id = ""
  instance_id = ""
}

```
Import

dcdb security_group_attachment can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_security_group_attachment.security_group_attachment securityGroupAttachment_id
```