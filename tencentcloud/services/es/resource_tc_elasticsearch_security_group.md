Provides a resource to create a elasticsearch security_group

Example Usage

```hcl
resource "tencentcloud_elasticsearch_security_group" "security_group" {
    instance_id        = "es-5wn36he6"
    security_group_ids = [
        "sg-mayqdlt1",
        "sg-po2q8cg7",
    ]
}
```

Import

elasticsearch security_group can be imported using the id, e.g.

```
terraform import tencentcloud_elasticsearch_security_group.security_group instance_id
```