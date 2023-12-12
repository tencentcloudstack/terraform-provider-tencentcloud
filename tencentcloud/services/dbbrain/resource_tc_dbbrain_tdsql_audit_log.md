Provides a resource to create a dbbrain tdsql_audit_log

Example Usage

```hcl
resource "tencentcloud_dbbrain_tdsql_audit_log" "my_log" {
  product = "dcdb"
  node_request_type = "dcdb"
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  filter {
		host = ["%%", "127.0.0.1"]
		user = ["tf_test", "mysql"]
  }
}
```