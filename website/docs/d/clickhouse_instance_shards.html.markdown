---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_instance_shards"
sidebar_current: "docs-tencentcloud-datasource-clickhouse_instance_shards"
description: |-
  Use this data source to query detailed information of clickhouse instance_shards
---

# tencentcloud_clickhouse_instance_shards

Use this data source to query detailed information of clickhouse instance_shards

## Example Usage

```hcl
data "tencentcloud_clickhouse_instance_shards" "instance_shards" {
  instance_id = "cdwch-datuhk3z"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Cluster instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_shards_list` - Instance shard information.


