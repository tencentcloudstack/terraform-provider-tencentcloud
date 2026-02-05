Provides a resource to create Security group.

Example Usage

Create a basic security group

```hcl
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg test"
}
```

Create a complete security group

```hcl
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg test"
  project_id  = 0

  tags = {
    "createdBy" = "Terraform"
  }
}
```

Import

Security group can be imported using the id, e.g.

```
terraform import tencentcloud_security_group.example sg-ey3wmiz1
```