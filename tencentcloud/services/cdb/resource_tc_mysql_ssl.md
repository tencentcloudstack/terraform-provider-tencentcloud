Provides a resource to create a MySQL SSL

Example Usage

For mysql instance SSL

```hcl
resource "tencentcloud_mysql_ssl" "example" {
  instance_id = "cdb-j5rprr8n"
  status      = "OFF"
}
```

For mysql RO group SSL

```hcl
resource "tencentcloud_mysql_ssl" "example" {
  ro_group_id = "cdbrg-k9a6gup3"
  status      = "ON"
}
```

Import

MySQL SSL can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_ssl.example cdb-j5rprr8n
```

Or

```
terraform import tencentcloud_mysql_ssl.example cdbrg-k9a6gup3
```
