Provides a resource to create a mysql rollback_stop

Example Usage

Revoke the ongoing rollback task of the instance

```hcl
resource "tencentcloud_mysql_rollback_stop" "example" {
  instance_id = "cdb-fitq5t9h"
}
```