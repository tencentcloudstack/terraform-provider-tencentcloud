---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_serverless_hbase_instances"
sidebar_current: "docs-tencentcloud-datasource-serverless_hbase_instances"
description: |-
  Use this data source to query detailed information of emr serverless_hbase_instances
---

# tencentcloud_serverless_hbase_instances

Use this data source to query detailed information of emr serverless_hbase_instances

## Example Usage

```hcl
data "tencentcloud_serverless_hbase_instances" "serverless_hbase_instances" {
  display_strategy = "clusterList"
}
```

## Argument Reference

The following arguments are supported:

* `display_strategy` - (Required, String) Cluster filtering policy. Value range:
	* clusterList: Query the list of clusters except the destroyed cluster;
	* monitorManage: Queries the list of clusters except those destroyed, being created, and failed to create.
* `asc` - (Optional, Int) Sort by OrderField in ascending or descending order. Value range:
	* 0: indicates the descending order;
	* 1: indicates the ascending order;
	The default value is 0.
* `filters` - (Optional, List) Custom query.
* `order_field` - (Optional, String) Sorting field. Value range:
	* clusterId: Sorting by instance ID;
	* addTime: sorted by instance creation time;
	* status: sorted by the status code of the instance.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field name.
* `values` - (Required, Set) Filter field value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - Cluster instance list.
  * `add_time` - Creation time.
  * `app_id` - User APP ID.
  * `cluster_id` - Cluster Instance String ID.
  * `cluster_name` - Cluster Instance name.
  * `id` - Cluster Instance Digital ID.
  * `pay_mode` - Cluster charging type. 0 means charging by volume, 1 means annual and monthly.
  * `region_id` - Region ID.
  * `status_desc` - State description.
  * `status` - Status code, please refer to the StatusDesc.
  * `subnet_id` - Primary Availability Subnet ID.
  * `tags` - List of tags.
  * `vpc_id` - Primary Availability Vpc ID.
  * `zone_id` - Primary Availability Zone ID.
  * `zone_settings` - Detailed configuration of the instance availability zone, including the availability zone name, VPC information, and the total number of nodes, where the total number of nodes must be greater than or equal to 3 and less than or equal to 50.
    * `node_num` - Number of nodes.
    * `vpc_settings` - Private network related information configuration. This parameter can be used to specify the ID of the private network, subnet ID, and other information.
      * `subnet_id` - Subnet ID.
      * `vpc_id` - VPC ID.
    * `zone` - The availability zone to which the instance belongs, such as ap-guangzhou-1.
  * `zone` - Primary Availability Zone Name.


