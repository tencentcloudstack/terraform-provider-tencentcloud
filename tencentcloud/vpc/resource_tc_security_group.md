Provides a resource to create security group.

Example Usage

Create a basic security group

```hcl
resource "tencentcloud_security_group" "example" {
  name        = "tf-example-sg"
  description = "sg test"
}
```

Create a complete security group

```hcl
resource "tencentcloud_security_group" "example" {
  name        = "tf-example-sg"
  description = "sg test"
  project_id  = 0

  tags = {
    "example" = "test"
  }
}
```

Import

Security group can be imported using the id, e.g.

```
  $ terraform import tencentcloud_security_group.sglab sg-ey3wmiz1
```