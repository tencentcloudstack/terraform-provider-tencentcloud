---
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

* `availability_zone` - (Required) The zone of the postgresql instance to query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of zones will be exported and its every element contains the following attributes:
  * `cpu` - The CPU number of the postgresql instance.
  * `id` - Id of the speccode of the postgresql instance. This parameter is used as `spec_code` for the creation of postgresql instance.
  * `memory` - Memory size(in MB).
  * `qps` - The QPS of the postgresql instance.
  * `storage_max` - The maximum volume size(in GB).
  * `storage_min` - The minimum volume size(in GB).
  * `version_name` - The version name of the postgresql instance.
  * `version` - The version of the postgresql instance.


