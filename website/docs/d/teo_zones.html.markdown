---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zones"
sidebar_current: "docs-tencentcloud-datasource-teo_zones"
description: |-
  Use this data source to query detailed information of teo zoneAvailablePlans
---

# tencentcloud_teo_zones

Use this data source to query detailed information of teo zoneAvailablePlans

## Example Usage

```hcl
data "tencentcloud_teo_zones" "teo_zones" {
  filters {
    name   = "zone-id"
    values = ["zone-39quuimqg8r6"]
  }

  filters {
    name   = "tag-key"
    values = ["createdBy"]
  }

  filters {
    name   = "tag-value"
    values = ["terraform"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Optional, String) Sort direction. If the field value is a number, sort by the numeric value. If the field value is text, sort by the ascill code. Values include: `asc`: From the smallest to largest; `desc`: From the largest to smallest. Default value: `desc`.
* `filters` - (Optional, List) Filter criteria. the maximum value of Filters.Values is 20. if this parameter is left empty, all site information authorized under the current appid will be returned. detailed filter criteria are as follows: zone-name: filter by site name; zone-id: filter by site id. the site id is in the format of zone-2noz78a8ev6k; status: filter by site status; tag-key: filter by tag key; tag-value: filter by tag value; alias-zone-name: filter by identical site identifier. when performing a fuzzy query, the fields that support filtering are named zone-name or alias-zone-name.
* `order` - (Optional, String) Sort the returned results according to this field. Values include: `type`: Connection mode; `area`: Acceleration region; `create-time`: Creation time; `zone-name`: Site name; `use-time`: Last used time; `active-status` Effective status. Default value: `create-time`.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Value of the filtered field.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zones` - Details of sites.


