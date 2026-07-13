---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_dms_database"
sidebar_current: "docs-tencentcloud-resource-dlc_dms_database"
description: |-
  Provides a resource to create a DMS database for Tencent Cloud DLC (Data Lake Compute).
---

# tencentcloud_dlc_dms_database

Provides a resource to create a DMS database for Tencent Cloud DLC (Data Lake Compute).

## Example Usage

```hcl
resource "tencentcloud_dlc_dms_database" "example" {
  name                       = "tf_example_dms_database"
  schema_name                = "DataLake"
  datasource_connection_name = "tf_example_connection"
  location                   = "cosn://tf-example-bucket-1300000000/data/"

  asset {
    name        = "tf_example_asset"
    catalog     = "DataLake"
    description = "example dms database asset."
    owner       = "root"

    params {
      key   = "param_key"
      value = "param_value"
    }

    biz_params {
      key   = "biz_key"
      value = "biz_value"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `datasource_connection_name` - (Required, String, ForceNew) Datasource connection name.
* `name` - (Required, String, ForceNew) Database name.
* `schema_name` - (Required, String, ForceNew) Schema name.
* `asset` - (Optional, List) Basic metadata object.
* `cascade` - (Optional, Bool) Whether to cascade delete when deleting the database.
* `delete_data` - (Optional, Bool) Whether to delete data when deleting the database.
* `location` - (Optional, String) Db storage path.

The `asset` object supports the following:

* `biz_params` - (Optional, List) Additional business attributes.
* `catalog` - (Optional, String) Data catalog.
* `description` - (Optional, String) Description.
* `name` - (Optional, String) Name.
* `owner_account` - (Optional, String) Object owner account.
* `owner` - (Optional, String) Object owner.
* `params` - (Optional, List) Additional attributes.
* `perm_values` - (Optional, List) Permissions.

The `biz_params` object of `asset` supports the following:

* `key` - (Required, String) Configured key value.
* `value` - (Optional, String) Configured value.

The `params` object of `asset` supports the following:

* `key` - (Required, String) Configured key value.
* `value` - (Optional, String) Configured value.

The `perm_values` object of `asset` supports the following:

* `key` - (Required, String) Configured key value.
* `value` - (Optional, String) Configured value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

DLC DMS database can be imported using the compound id `name#schema_name#datasource_connection_name`, e.g.

```
terraform import tencentcloud_dlc_dms_database.example tf_example_dms_database#DataLake#tf_example_connection
```

