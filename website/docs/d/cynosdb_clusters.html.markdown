---
subcategory: "CynosDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_clusters"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_clusters"
description: |-
  Use this data source to query detailed information of Cynosdb clusters.
---

# tencentcloud_cynosdb_clusters

Use this data source to query detailed information of Cynosdb clusters.

## Example Usage

```hcl
data "tencentcloud_cynosdb_clusters" "foo" {
  cluster_id   = "cynosdbmysql-dzj5l8gz"
  project_id   = 0
  db_type      = "MYSQL"
  cluster_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional, String) ID of the cluster to be queried.
* `cluster_name` - (Optional, String) Name of the cluster to be queried.
* `db_type` - (Optional, String) Type of CynosDB, and available values include `MYSQL`, `POSTGRESQL`.
* `project_id` - (Optional, Int) ID of the project to be queried.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_list` - A list of clusters. Each element contains the following attributes:
  * `auto_renew_flag` - Auto renew flag. Valid values are `0`(MANUAL_RENEW), `1`(AUTO_RENEW). Only works for PREPAID cluster.
  * `available_zone` - The available zone of the CynosDB Cluster.
  * `charge_type` - The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`.
  * `cluster_id` - ID of CynosDB cluster.
  * `cluster_limit` - Storage limit of CynosDB cluster instance, unit in GB.
  * `cluster_name` - Name of CynosDB cluster.
  * `cluster_status` - Status of the Cynosdb cluster.
  * `create_time` - Creation time of the CynosDB cluster.
  * `db_type` - Type of CynosDB, and available values include `MYSQL`.
  * `db_version` - Version of CynosDB, which is related to `db_type`. For `MYSQL`, available value is `5.7`.
  * `port` - Port of CynosDB cluster.
  * `project_id` - ID of the project.
  * `subnet_id` - ID of the subnet within this VPC.
  * `vpc_id` - ID of the VPC.


