---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_spec"
sidebar_current: "docs-tencentcloud-datasource-clickhouse_spec"
description: |-
  Use this data source to query detailed information of clickhouse spec
---

# tencentcloud_clickhouse_spec

Use this data source to query detailed information of clickhouse spec

## Example Usage

```hcl
data "tencentcloud_clickhouse_spec" "spec" {
  zone       = "ap-guangzhou-7"
  pay_mode   = "PREPAID"
  is_elastic = false
}
```

## Argument Reference

The following arguments are supported:

* `zone` - (Required, String) Regional information.
* `is_elastic` - (Optional, Bool) Is it elastic.
* `pay_mode` - (Optional, String) Billing type, PREPAID means annual and monthly subscription, POSTPAID_BY_HOUR means pay-as-you-go billing.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `attach_cbs_spec` - Cloud disk list.
  * `disk_count` - Number of disks.
  * `disk_desc` - Disk type description.
  * `disk_type` - Disk type.
  * `max_disk_size` - Maximum disk size, unit G.
  * `min_disk_size` - Minimum disk size, unit G.
* `common_spec` - Zookeeper node specification description.
  * `available` - Whether it is available, false means sold out.
  * `compute_spec_desc` - Specification description information.
  * `cpu` - Number of cpu cores.
  * `data_disk` - Data disk description information.
    * `disk_count` - Number of disks.
    * `disk_desc` - Disk type description.
    * `disk_type` - Disk type.
    * `max_disk_size` - Maximum disk size, unit G.
    * `min_disk_size` - Minimum disk size, unit G.
  * `display_name` - Specification name.
  * `instance_quota` - Inventory.
  * `max_node_size` - Maximum number of nodes limit.
  * `mem` - Memory size, unit G.
  * `name` - Specification name.
  * `system_disk` - System disk description information.
    * `disk_count` - Number of disks.
    * `disk_desc` - Disk type description.
    * `disk_type` - Disk type.
    * `max_disk_size` - Maximum disk size, unit G.
    * `min_disk_size` - Minimum disk size, unit G.
  * `type` - Classification tags, STANDARD/BIGDATA/HIGHIO respectively represent standard/big data/high IO.
* `data_spec` - Data node specification description.
  * `available` - Whether it is available, false means sold out.
  * `compute_spec_desc` - Specification description information.
  * `cpu` - Number of cpu cores.
  * `data_disk` - Data disk description information.
    * `disk_count` - Number of disks.
    * `disk_desc` - Disk type description.
    * `disk_type` - Disk type.
    * `max_disk_size` - Maximum disk size, unit G.
    * `min_disk_size` - Minimum disk size, unit G.
  * `display_name` - Specification name.
  * `instance_quota` - Inventory.
  * `max_node_size` - Maximum number of nodes limit.
  * `mem` - Memory size, unit G.
  * `name` - Specification name.
  * `system_disk` - System disk description information.
    * `disk_count` - Number of disks.
    * `disk_desc` - Disk type description.
    * `disk_type` - Disk type.
    * `max_disk_size` - Maximum disk size, unit G.
    * `min_disk_size` - Minimum disk size, unit G.
  * `type` - Classification tags, STANDARD/BIGDATA/HIGHIO respectively represent standard/big data/high IO.


