---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_db_instance_ssl_config"
sidebar_current: "docs-tencentcloud-resource-sqlserver_db_instance_ssl_config"
description: |-
  Provides a resource to create a SQL Server db instance ssl config
---

# tencentcloud_sqlserver_db_instance_ssl_config

Provides a resource to create a SQL Server db instance ssl config

## Example Usage

### Enable SSL encryption for SQL Server instance

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  type        = "enable"
}
```

### Disable SSL encryption for SQL Server instance

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  type        = "disable"
}
```

### Renew SSL certificate for SQL Server instance

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  type        = "renew"
}
```

### Enable SSL encryption with KMS protection

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  type        = "enable"
  is_kms      = 1
  key_id      = "your-cmk-key-id"
  key_region  = "ap-guangzhou"
}
```

### Enable SSL encryption during maintenance window

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  type        = "enable"
  wait_switch = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) SQL Server instance ID.
* `type` - (Required, String) SSL operation type. Valid values: enable, disable, renew.
* `is_kms` - (Optional, Int) Whether to enable KMS encryption protection. 0: no, 1: yes. Default is 0.
* `key_id` - (Optional, String) KMS CMK key ID, required when IsKMS is 1.
* `key_region` - (Optional, String) CMK region, required when IsKMS is 1.
* `wait_switch` - (Optional, Int) Execution timing. 0: execute immediately, 1: execute during maintenance window. Default is 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `encryption` - SSL encryption status. Valid values: enable, disable, enable_doing, disable_doing, renew_doing, wait_doing.
* `ssl_validity_period` - SSL certificate validity period, format: YYYY-MM-DD HH:MM:SS.
* `ssl_validity` - SSL certificate validity. 0: invalid, 1: valid.


## Import

SQL Server db instance ssl config can be imported using the instance_id, e.g.

```
terraform import tencentcloud_sqlserver_db_instance_ssl_config.example mssql-gy1lc54f
```

