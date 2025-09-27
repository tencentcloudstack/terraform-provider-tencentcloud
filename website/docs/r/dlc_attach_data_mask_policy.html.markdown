---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_attach_data_mask_policy"
sidebar_current: "docs-tencentcloud-resource-dlc_attach_data_mask_policy"
description: |-
  Provides a resource to create a DLC attach data mask policy
---

# tencentcloud_dlc_attach_data_mask_policy

Provides a resource to create a DLC attach data mask policy

## Example Usage

```hcl
resource "tencentcloud_dlc_data_mask_strategy" "example" {
  strategy {
    strategy_name = "tf-example"
    strategy_desc = "description."
    groups {
      work_group_id = 70220
      strategy_type = "MASK"
    }
  }
}

resource "tencentcloud_dlc_attach_data_mask_policy" "example" {
  data_mask_strategy_policy_set {
    policy_info {
      database = "tf-example"
      catalog  = "DataLakeCatalog"
      table    = "tf-example"
      column   = "id"
    }

    data_mask_strategy_id = tencentcloud_dlc_data_mask_strategy.example.id
    column_type           = "string"
  }
}
```

## Argument Reference

The following arguments are supported:

* `data_mask_strategy_policy_set` - (Optional, List, ForceNew) A collection of data masking policy permission objects to be bound.

The `data_mask_strategy_policy_set` object supports the following:

* `column_type` - (Optional, String, ForceNew) The type of the bound field.
* `data_mask_strategy_id` - (Optional, String, ForceNew) The ID of the data masking strategy.
* `policy_info` - (Optional, List, ForceNew) Data masking permission object.

The `policy_info` object of `data_mask_strategy_policy_set` supports the following:

* `catalog` - (Required, String, ForceNew) The name of the data source to be authorized. For administrator-level authorization, only * is allowed (representing all resources at this level). For data source-level and database-level authorization, only COSDataCatalog or * is allowed. For table-level authorization, custom data sources can be specified. Defaults to DataLakeCatalog if not specified. Note: For custom data sources, DLC can only manage a subset of permissions provided by the user during data source integration.
* `database` - (Required, String, ForceNew) The name of the database to be authorized. Use * to represent all databases under the current Catalog. For administrator-level authorization, only * is allowed. For data connection-level authorization, leave it empty. For other types, specify the database name.
* `table` - (Required, String, ForceNew) The name of the table to be authorized. Use * to represent all tables under the current Database. For administrator-level authorization, only * is allowed. For data connection-level and database-level authorization, leave it empty. For other types, specify the table name.
* `column` - (Optional, String, ForceNew) The name of the column to be authorized. Use * to represent all columns. For administrator-level authorization, only * is allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



