Provides a resource to create a mysql ssl

Example Usage

```hcl
resource "tencentcloud_mysql_ssl" "ssl" {
  instance_id = "cdb-j5rprr8n"
  status      = "OFF"
}
```

Import

mysql ssl can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_ssl.ssl instanceId
```