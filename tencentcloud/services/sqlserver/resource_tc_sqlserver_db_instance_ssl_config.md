Provides a resource to create a SQL Server db instance ssl config

Example Usage

Enable SSL encryption for SQL Server instance

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  type        = "enable"
}
```

Disable SSL encryption for SQL Server instance

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  type        = "disable"
}
```

Renew SSL certificate for SQL Server instance

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  type        = "renew"
}
```

Enable SSL encryption with KMS protection

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  type        = "enable"
  is_kms      = 1
  key_id      = "your-cmk-key-id"
  key_region  = "ap-guangzhou"
}
```

Enable SSL encryption during maintenance window

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  type        = "enable"
  wait_switch = 1
}
```

Import

SQL Server db instance ssl config can be imported using the instance_id, e.g.

```
terraform import tencentcloud_sqlserver_db_instance_ssl_config.example mssql-gy1lc54f
```
