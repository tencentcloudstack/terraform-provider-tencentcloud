---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_zones"
sidebar_current: "docs-tencentcloud-datasource-tcaplus_zones"
description: |-
  Use this data source to query tcaplus zones
---

# tencentcloud_tcaplus_zones

Use this data source to query tcaplus zones

## Example Usage

```hcl
data "tencentcloud_tcaplus_zones" "null" {
  app_id = "19162256624"
}
data "tencentcloud_tcaplus_zones" "id" {
  app_id  = "19162256624"
  zone_id = "19162256624:1"
}
data "tencentcloud_tcaplus_zones" "name" {
  app_id    = "19162256624"
  zone_name = "test"
}
data "tencentcloud_tcaplus_zones" "all" {
  app_id    = "19162256624"
  zone_id   = "19162256624:1"
  zone_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) Id of the tcapplus application to be query.
* `result_output_file` - (Optional) Used to save results.
* `zone_id` - (Optional) Zone id to be query.
* `zone_name` - (Optional) Zone name to be query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of tcaplus zones. Each element contains the following attributes.
  * `create_time` - create time of the tcapplus application.
  * `table_count` - Number of tables.
  * `total_size` - The total storage(MB).
  * `zone_id` - Id of the tcapplus zone.
  * `zone_name` - Name of the tcapplus zone.


