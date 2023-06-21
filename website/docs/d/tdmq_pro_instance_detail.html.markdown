---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_pro_instance_detail"
sidebar_current: "docs-tencentcloud-datasource-tdmq_pro_instance_detail"
description: |-
  Use this data source to query detailed information of tdmq pro_instance_detail
---

# tencentcloud_tdmq_pro_instance_detail

Use this data source to query detailed information of tdmq pro_instance_detail

## Example Usage

```hcl
data "tencentcloud_tdmq_pro_instance_detail" "pro_instance_detail" {
  cluster_id = "pulsar-9n95ax58b9vn"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster Id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_info` - Cluster information.
  * `can_edit_route` - Can the route be modifiedNote: This field may return null, indicating that no valid value can be obtained.
  * `cluster_id` - Cluster Id.
  * `cluster_name` - Cluster name.
  * `create_time` - Creation time.
  * `max_storage` - Maximum storage capacity, unit: MB.
  * `node_distribution` - Node distributionNote: This field may return null, indicating that no valid value can be obtained.
    * `node_count` - Number of nodes.
    * `zone_id` - Availability zone ID.
    * `zone_name` - Availability zone.
  * `remark` - Descriptive information.
  * `status` - Cluster status, 0: creating, 1: normal, 2: isolated.
  * `version` - cluster version.
* `cluster_spec_info` - Cluster specification informationNote: This field may return null, indicating that no valid value can be obtained.
  * `max_band_width` - peak bandwidth. Unit: mbps.
  * `max_namespaces` - Maximum number of namespaces.
  * `max_topics` - Maximum number of topic partitions.
  * `max_tps` - peak tps.
  * `scalable_tps` - Elastic TPS outside specificationNote: This field may return null, indicating that no valid value can be obtained.
  * `spec_name` - Cluster specification name.
* `network_access_point_infos` - Cluster network access point informationNote: This field may return null, indicating that no valid value can be obtained.
  * `endpoint` - access address.
  * `instance_id` - instance id.
  * `route_type` - Access point type: 0: support network access point 1: VPC access point 2: public network access point.
  * `subnet_id` - Subnet id, support network and public network access point, this field is emptyNote: This field may return null, indicating that no valid value can be obtained.
  * `vpc_id` - The id of the vpc, the supporting network and the access point of the public network, this field is emptyNote: This field may return null, indicating that no valid value can be obtained.


