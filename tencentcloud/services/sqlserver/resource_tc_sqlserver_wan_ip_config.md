Provides a resource to create a sqlserver wan ip config

Example Usage

Open wan ip for SQL instance

```hcl
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  enable_wan_ip = true
}
```

Close wan ip for SQL instance

```hcl
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  enable_wan_ip = false
}
```

Open wan ip for SQL read only group

```hcl
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  ro_group_id   = "mssqlrg-hyxotm31"
  enable_wan_ip = true
}
```

Close wan ip for SQL read only group

```hcl
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  ro_group_id   = "mssqlrg-hyxotm31"
  enable_wan_ip = false
}
```
