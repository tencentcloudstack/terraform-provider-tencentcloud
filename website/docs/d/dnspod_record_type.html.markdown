---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_record_type"
sidebar_current: "docs-tencentcloud-datasource-dnspod_record_type"
description: |-
  Use this data source to query detailed information of dnspod record_type
---

# tencentcloud_dnspod_record_type

Use this data source to query detailed information of dnspod record_type

## Example Usage

```hcl
data "tencentcloud_dnspod_record_type" "record_type" {
  domain_grade = "DP_FREE"
}
```

## Argument Reference

The following arguments are supported:

* `domain_grade` - (Required, String) Domain level. + Old packages: D_FREE, D_PLUS, D_EXTRA, D_EXPERT, D_ULTRA correspond to free package, personal luxury, enterprise 1, enterprise 2, enterprise 3. + New packages: DP_FREE, DP_PLUS, DP_EXTRA, DP_EXPERT, DP_ULTRA correspond to new free, personal professional, enterprise basic, enterprise standard, enterprise flagship.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type_list` - Record type list.


