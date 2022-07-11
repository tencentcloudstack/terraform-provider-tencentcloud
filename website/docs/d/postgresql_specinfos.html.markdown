---
subcategory: "PostgreSQL"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_specinfos"
sidebar_current: "docs-tencentcloud-datasource-postgresql_specinfos"
description: |-
  Use this data source to get the available product configs of the postgresql instance.
---

# tencentcloud_postgresql_specinfos

Use this data source to get the available product configs of the postgresql instance.

## Example Usage

```hcl
data "tencentcloud_postgresql_specinfos" "foo" {
  availability_zone = "ap-shanghai-2"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, String) The zone of the postgresql instance to query.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of zones will be exported and its every element contains the following attributes:
  * `cpu` - The CPU number of the postgresql instance.
  * `engine_version_name` - Version name of the postgresql database engine.
  * `engine_version` - Version of the postgresql database engine.
  * `id` - ID of the postgresql instance speccode.
  * `memory` - Memory size(in GB).
  * `qps` - The QPS of the postgresql instance.
  * `storage_max` - The maximum volume size(in GB).
  * `storage_min` - The minimum volume size(in GB).


