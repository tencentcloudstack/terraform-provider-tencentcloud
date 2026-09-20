---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_tc_lake_meta_instance"
sidebar_current: "docs-tencentcloud-datasource-dlc_tc_lake_meta_instance"
description: |-
  Use this data source to query the activation status of DLC (Data Lake Compute) TCLake meta instance.
---

# tencentcloud_dlc_tc_lake_meta_instance

Use this data source to query the activation status of DLC (Data Lake Compute) TCLake meta instance.

## Example Usage

```hcl
data "tencentcloud_dlc_tc_lake_meta_instance" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - TCLake activation status. Enumerated value: `Running` means activated successfully.


