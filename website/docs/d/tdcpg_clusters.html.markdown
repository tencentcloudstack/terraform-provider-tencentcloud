---
subcategory: "TDSQL-C for PostgreSQL(TDCPG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdcpg_clusters"
sidebar_current: "docs-tencentcloud-datasource-tdcpg_clusters"
description: |-
  Use this data source to query detailed information of tdcpg clusters.
---

# tencentcloud_tdcpg_clusters

Use this data source to query detailed information of tdcpg clusters.

~> **NOTE:** This data source is still in internal testing. To experience its functions, you need to apply for a whitelist from Tencent Cloud.

## Example Usage

```hcl
data "tencentcloud_tdcpg_clusters" "clusters" {
  cluster_id   = ""
  cluster_name = ""
  status       = ""
  pay_mode     = ""
  project_id   = ""
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional, String) cluster id.
* `cluster_name` - (Optional, String) cluster name.
* `pay_mode` - (Optional, String) pay mode.
* `project_id` - (Optional, Int) project id, default to 0, means default project.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, String) cluster status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - cluster list.
  * `auto_renew_flag` - auto renew flag.
  * `cluster_id` - cluster id.
  * `cluster_name` - cluster name.
  * `create_time` - create time.
  * `db_charset` - db charset.
  * `db_kernel_version` - db kernel version.
  * `db_major_version` - db major version.
  * `db_version` - db version.
  * `endpoint_set` - endpoint set.
    * `cluster_id` - cluster id.
    * `endpoint_id` - endpoint id.
    * `endpoint_name` - endpoint name.
    * `endpoint_type` - endpoint type.
    * `private_ip` - private ip.
    * `private_port` - private port.
    * `subnet_id` - subnet id.
    * `vpc_id` - vpc id.
    * `wan_domain` - wan domain.
    * `wan_ip` - wan ip.
    * `wan_port` - wan port.
  * `instance_count` - instance count.
  * `pay_mode` - pay mode.
  * `pay_period_end_time` - pay period expired time.
  * `project_id` - project id.
  * `region` - region.
  * `status_desc` - status description.
  * `status` - status.
  * `storage_limit` - storage limit, unit is GB.
  * `storage_pay_mode` - storage pay mode, optional value is PREPAID or POSTPAID_BY_HOUR.
  * `storage_used` - storage used, unit is GB.
  * `zone` - zone.


