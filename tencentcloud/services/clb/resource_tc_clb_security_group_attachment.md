Provides a resource to create a clb security_group_attachment

Example Usage

```hcl
resource "tencentcloud_clb_security_group_attachment" "security_group_attachment" {
  security_group = "sg-ijato2x1"
  load_balancer_ids = ["lb-5dnrkgry"]
}
```

Import

clb security_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_clb_security_group_attachment.security_group_attachment security_group_id#clb_id
```