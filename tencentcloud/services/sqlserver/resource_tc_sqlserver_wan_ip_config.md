Provides a resource to create a sqlserver wan ip config

Example Usage

Open/Close wan ip for SQL instance

```hcl
# open
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  enable_wan_ip = true
}

# close
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  enable_wan_ip = false
}
```

Open/Close wan ip for SQL read only group

```hcl
# open
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  ro_group_id   = "mssqlrg-hyxotm31"
  enable_wan_ip = true
}

# close
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  ro_group_id   = "mssqlrg-hyxotm31"
  enable_wan_ip = false
}
```
