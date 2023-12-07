Provides a resource to create a dcdb db_parameters

Example Usage

```hcl
resource "tencentcloud_dcdb_db_parameters" "db_parameters" {
  instance_id = "%s"
  params {
	param = "max_connections"
	value = "9999"
  }
}
```

Import

dcdb db_parameters can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_db_parameters.db_parameters instanceId#paramName
```