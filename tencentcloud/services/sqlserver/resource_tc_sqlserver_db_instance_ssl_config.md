Provides a resource to create a SQL Server db instance ssl config

Example Usage

Enable SSL encryption for SQL Server instance

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  encryption  = "enable"
}
```

Disable SSL encryption for SQL Server instance

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  encryption  = "disable"
}
```

Enable SSL encryption with KMS protection

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  encryption  = "enable"
  is_kms      = 1
  cmk_id      = "your-cmk-key-id"
  cmk_region  = "ap-guangzhou"
}
```

Import

SQL Server db instance ssl config can be imported using the instance_id, e.g.

```
terraform import tencentcloud_sqlserver_db_instance_ssl_config.example mssql-gy1lc54f
```
