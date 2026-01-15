---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_quality_rule_templates"
sidebar_current: "docs-tencentcloud-datasource-wedata_quality_rule_templates"
description: |-
  Use this data source to query detailed information of WeData quality rule templates.
---

# tencentcloud_wedata_quality_rule_templates

Use this data source to query detailed information of WeData quality rule templates.

## Example Usage

```hcl
data "tencentcloud_wedata_quality_rule_templates" "example" {
  project_id = "your_project_id"

  order_fields {
    name      = "CitationCount"
    direction = "DESC"
  }

  filters {
    name   = "Type"
    values = ["1", "2"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Workspace ID.
* `filters` - (Optional, List) General filter conditions 1. `Id` Description: Template ID, Value: Unique identifier of the template; 2. `Keyword` Description: Keyword search, supports fuzzy search of template names, Value: String; 3. `Type` Description: Template type, Value: `1` - System template; `2` - Custom template; supports multiple values (OR relationship); 4. `QualityDim` Description: Quality detection dimension, Value: `1` - Accuracy; `2` - Uniqueness; `3` - Completeness; `4` - Consistency; `5` - Timeliness; `6` - Validity; supports multiple values (OR relationship); 5. `SourceObjectType` Description: Source data object type applicable to the rule, Value: `1` - Constant; `2` - Offline table level; `3` - Offline field level; `4` - Database level; supports multiple values (OR relationship); 6. `SourceEngineTypes` Description: Source data engine type applicable to the template, Value: `1` - MySQL; `2` - Hive; `4` - Spark; `8` - Livy; `16` - DLC; `32` - Gbase; `64` - TCHouse-P; `128` - Doris; `256` - TCHouse-D; `512` - EMR_StarRocks; `1024` - TCHouse-X; supports multiple values (OR relationship).
* `order_fields` - (Optional, List) General sorting, supported sorting fields: `CitationCount` - Sort by citation count; `UpdateTime` - Sort by update time. Sort direction: `1` - Ascending (ASC); `2` - Descending (DESC).
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) Filter field name.
* `values` - (Optional, Set) Filter value list.

The `order_fields` object supports the following:

* `direction` - (Required, String) Sort direction: ASC|DESC.
* `name` - (Required, String) Sort field name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Result.


