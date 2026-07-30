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
  encryption  = "enable"
}
```

### Disable SSL encryption for SQL Server instance

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  encryption  = "disable"
}
```

### Enable SSL encryption with KMS protection

```hcl
resource "tencentcloud_sqlserver_db_instance_ssl_config" "example" {
  instance_id = "mssql-gy1lc54f"
  encryption  = "enable"
  is_kms      = 1
  cmk_id      = "your-cmk-key-id"
  cmk_region  = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `encryption` - (Required, String) SSL encryption desired state. Valid values: `enable`, `disable`.
* `instance_id` - (Required, String, ForceNew) SQL Server instance ID.
* `cmk_id` - (Optional, String) KMS CMK key ID, required when is_kms is 1.
* `cmk_region` - (Optional, String) CMK region, required when is_kms is 1.
* `is_kms` - (Optional, Int) Whether to enable KMS encryption protection. 0: no, 1: yes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `ssl_validity_period` - SSL certificate validity period, format: YYYY-MM-DD HH:MM:SS.
* `ssl_validity` - SSL certificate validity. 0: invalid, 1: valid.


## Import

SQL Server db instance ssl config can be imported using the instance_id, e.g.

```
terraform import tencentcloud_sqlserver_db_instance_ssl_config.example mssql-gy1lc54f
```

