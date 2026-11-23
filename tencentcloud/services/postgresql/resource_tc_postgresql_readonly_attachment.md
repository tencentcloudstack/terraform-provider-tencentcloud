Use this resource to create postgresql readonly attachment.

Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_attachment" "example" {
  db_instance_id     = "pgro-k95qzetn"
  read_only_group_id = "pgrogrp-test"
}
```