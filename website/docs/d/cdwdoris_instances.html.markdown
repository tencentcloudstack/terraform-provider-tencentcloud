---
subcategory: "CdwDoris"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwdoris_instances"
sidebar_current: "docs-tencentcloud-datasource-cdwdoris_instances"
description: |-
  Use this data source to query detailed information of CDWDoris instances
---

# tencentcloud_cdwdoris_instances

Use this data source to query detailed information of CDWDoris instances

## Example Usage

### Query all cdwdoris instances

```hcl
data "tencentcloud_cdwdoris_instances" "example" {}
```

### Query cdwdoris instances by filter

```hcl
# by instance Id
data "tencentcloud_cdwdoris_instances" "example" {
  search_instance_id = "cdwdoris-rhbflamd"
}

# by instance name
data "tencentcloud_cdwdoris_instances" "example" {
  search_instance_name = "tf-example"
}

# by instance tags
data "tencentcloud_cdwdoris_instances" "example" {
  search_tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
    all_value = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `search_instance_id` - (Optional, String) The name of the cluster ID for the search.
* `search_instance_name` - (Optional, String) The cluster name for the search.
* `search_tags` - (Optional, List) Search tag list.

The `search_tags` object supports the following:

* `all_value` - (Optional, Int) 1 means only the tag key is entered without a value, and 0 means both the key and the value are entered.
* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances_list` - Quantities of instances array.
  * `access_info` - Access address. Example: 10.0.0.1:9000. Note: This field may return null, indicating that no valid values can be obtained.
  * `bind_s_gs` - Bound security group information. Note: This field may return null, indicating that no valid values can be obtained.
  * `build_version` - Minor versions. Note: This field may return null, indicating that no valid values can be obtained.
  * `can_attach_cbs` - cbs. Note: This field may return null, indicating that no valid values can be obtained.
  * `case_sensitive` - Whether the table name is case sensitive, 0 refers to sensitive, 1 refers to insensitive, compared in lowercase; 2 refers to insensitive, and the table name is changed to lowercase for storage.. Note: This field may return null, indicating that no valid values can be obtained.
  * `characteristic` - Page features, used to block some page entrances on the front end.. Note: This field may return null, indicating that no valid values can be obtained.
  * `cls_log_set_id` - Logset ID. Note: This field may return null, indicating that no valid values can be obtained.
  * `cls_topic_id` - Log Topic ID. Note: This field may return null, indicating that no valid values can be obtained.
  * `components` - Component Information. Note: The return type here is map[string]struct, not the string type displayed. You can refer to Sample Value to parse the data. Note: This field may return null, indicating that no valid values can be obtained.
  * `cool_down_bucket` - COS buckets are used for hot and cold stratification. Note: This field may return null, indicating that no valid values can be obtained.
  * `core_summary` - Zookeeper node description information. Note: This field may return null, indicating that no valid values can be obtained.
    * `attach_cbs_spec` - Information of mounted cloud disks. Note: This field may return null, indicating that no valid values can be obtained.
      * `disk_count` - Total number of disks.
      * `disk_desc` - Description.
      * `disk_size` - Disk capacity, in GB.
      * `disk_type` - Node disk type, such as CLOUD_SSD, CLOUD_PREMIUM.
    * `core` - Number of CPU cores, in counts.
    * `disk_count` - Disk size. Note: This field may return null, indicating that no valid values can be obtained.
    * `disk_desc` - Disk description.
    * `disk_type` - Disk type.
    * `disk` - Disk size, in GB.
    * `encrypt` - Whether it is encrypted.. Note: This field may return null, indicating that no valid values can be obtained.
    * `max_disk_size` - Maximum disk. Note: This field may return null, indicating that no valid values can be obtained.
    * `memory` - Memory size, in GB.
    * `node_size` - Number of nodes.
    * `spec_core` - Specified cores. Note: This field may return null, indicating that no valid values can be obtained.
    * `spec_memory` - Specified memory. Note: This field may return null, indicating that no valid values can be obtained.
    * `spec` - Model, such as S1.
    * `sub_product_type` - Sub-product name. Note: This field may return null, indicating that no valid values can be obtained.
  * `cos_bucket_name` - COS bucket. Note: This field may return null, indicating that no valid values can be obtained.
  * `cos_move_factor` - Cold and hot stratification coefficient. Note: This field may return null, indicating that no valid values can be obtained.
  * `create_time` - Creation time. Note: This field may return null, indicating that no valid values can be obtained.
  * `eip` - Elastic network interface address. Note: This field may return null, indicating that no valid values can be obtained.
  * `enable_cool_down` - Whether to enable hot and cold stratification. 0 refers to disabled, and 1 refers to enabled.. Note: This field may return null, indicating that no valid values can be obtained.
  * `enable_multi_zones` - Whether it is a multi-AZ.. Note: This field may return null, indicating that no valid values can be obtained.
  * `enable_xml_config` - Whether to support XML configuration management. Note: This field may return null, indicating that no valid values can be obtained.
  * `expire_time` - Expiration time. Note: This field may return null, indicating that no valid values can be obtained.
  * `flow_msg` - Error process description information. Note: This field may return null, indicating that no valid values can be obtained.
  * `grace_shutdown_wait_seconds` - The timeout time for the graceful restart of the kernel. If it is -1, it means it is not set.. Note: This field may return null, indicating that no valid values can be obtained.
  * `ha_type` - High availability type: 0: non-high availability. 1: read high availability. 2: read-write high availability. Note: This field may return null, indicating that no valid values can be obtained.
  * `ha` - High availability, being true or false. Note: This field may return null, indicating that no valid values can be obtained.
  * `has_cls_topic` - Whether to enable logs. Note: This field may return null, indicating that no valid values can be obtained.
  * `id` - Record ID, in numerical type. Note: This field may return null, indicating that no valid values can be obtained.
  * `instance_id` - Cluster instance ID, `cdw-xxxx` string type. Note: This field may return null, indicating that no valid values can be obtained.
  * `instance_name` - Cluster instance name. Note: This field may return null, indicating that no valid values can be obtained.
  * `is_white_s_gs` - Whether users can bind security groups.. Note: This field may return null, indicating that no valid values can be obtained.
  * `kind` - external/local/yunti. Note: This field may return null, indicating that no valid values can be obtained.
  * `master_summary` - Data node description information. Note: This field may return null, indicating that no valid values can be obtained.
    * `attach_cbs_spec` - Information of mounted cloud disks. Note: This field may return null, indicating that no valid values can be obtained.
      * `disk_count` - Total number of disks.
      * `disk_desc` - Description.
      * `disk_size` - Disk capacity, in GB.
      * `disk_type` - Node disk type, such as CLOUD_SSD, CLOUD_PREMIUM.
    * `core` - Number of CPU cores, in counts.
    * `disk_count` - Disk size. Note: This field may return null, indicating that no valid values can be obtained.
    * `disk_desc` - Disk description.
    * `disk_type` - Disk type.
    * `disk` - Disk size, in GB.
    * `encrypt` - Whether it is encrypted.. Note: This field may return null, indicating that no valid values can be obtained.
    * `max_disk_size` - Maximum disk. Note: This field may return null, indicating that no valid values can be obtained.
    * `memory` - Memory size, in GB.
    * `node_size` - Number of nodes.
    * `spec_core` - Specified cores. Note: This field may return null, indicating that no valid values can be obtained.
    * `spec_memory` - Specified memory. Note: This field may return null, indicating that no valid values can be obtained.
    * `spec` - Model, such as S1.
    * `sub_product_type` - Sub-product name. Note: This field may return null, indicating that no valid values can be obtained.
  * `monitor` - Monitoring Information. Note: This field may return null, indicating that no valid values can be obtained.
  * `pay_mode` - Payment type: hour and prepay. Note: This field may return null, indicating that no valid values can be obtained.
  * `region_desc` - Region. Note: This field may return null, indicating that no valid values can be obtained.
  * `region_id` - Region ID, indicating the region. Note: This field may return null, indicating that no valid values can be obtained.
  * `region` - Region, ap-guangzhou. Note: This field may return null, indicating that no valid values can be obtained.
  * `renew_flag` - Automatic renewal marker. Note: This field may return null, indicating that no valid values can be obtained.
  * `restart_timeout` - Timeout period, in seconds. Note: This field may return null, indicating that no valid values can be obtained.
  * `status_desc` - Status description, such as `running`. Note: This field may return null, indicating that no valid values can be obtained.
  * `status` - Status,. Init is being created. Serving is running. Deleted indicates the cluster has been terminated. Deleting indicates the cluster is being terminated. Modify indicates the cluster is being changed. Note: This field may return null, indicating that no valid values can be obtained.
  * `subnet_id` - Subnet name. Note: This field may return null, indicating that no valid values can be obtained.
  * `tags` - Tag list. Note: This field may return null, indicating that no valid values can be obtained.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `user_network_infos` - User availability zone and subnet information. Note: This field may return null, indicating that no valid values can be obtained.
  * `version` - Version. Note: This field may return null, indicating that no valid values can be obtained.
  * `vpc_id` - VPC name. Note: This field may return null, indicating that no valid values can be obtained.
  * `zone_desc` - Note about availability zone, such as Guangzhou Zone 2. Note: This field may return null, indicating that no valid values can be obtained.
  * `zone` - Availability zone, ap-guangzhou-3. Note: This field may return null, indicating that no valid values can be obtained.


