---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_standard_engine_resource_group_config_information"
sidebar_current: "docs-tencentcloud-datasource-dlc_standard_engine_resource_group_config_information"
description: |-
  Use this data source to query detailed information of DLC standard engine resource group config information
---

# tencentcloud_dlc_standard_engine_resource_group_config_information

Use this data source to query detailed information of DLC standard engine resource group config information

## Example Usage

```hcl
data "tencentcloud_dlc_standard_engine_resource_group_config_information" "example" {
  sort_by = "create-time"
  sorting = "desc"
  filters {
    name = "engine-id"
    values = [
      "DataEngine-5plqp7q7"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions are optional, engine-resource-group-id or engine-id.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) Sort Field.
* `sorting` - (Optional, String) Ascending or descending.

The `filters` object supports the following:

* `name` - (Required, String) Attribute name. If there are multiple filters, the relationship between the filters is a logical OR relationship.
* `values` - (Required, Set) Attribute value, if there are multiple Values in the same Filter, the relationship between the Values under the same Filter is a logical OR relationship.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `standard_engine_resource_group_config_infos` - Standard engine resource group, configuration related information.


