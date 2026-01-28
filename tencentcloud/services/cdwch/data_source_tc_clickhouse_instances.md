Provides a list of ClickHouse (TCHouse-C) instances.

## Example Usage

### Query all instances

```hcl
data "tencentcloud_clickhouse_instances" "all" {
}
```

### Query by instance ID

```hcl
data "tencentcloud_clickhouse_instances" "by_id" {
  instance_id = "cdwch-xxxxxx"
}
```

### Query by instance name

```hcl
data "tencentcloud_clickhouse_instances" "by_name" {
  instance_name = "my-clickhouse-cluster"
}
```

### Query by tags

```hcl
data "tencentcloud_clickhouse_instances" "by_tags" {
  tags = {
    env = "production"
    app = "analytics"
  }
}
```

### Query with multiple filters

```hcl
data "tencentcloud_clickhouse_instances" "filtered" {
  instance_name = "test"
  tags = {
    env = "test"
  }
  is_simple = true
  result_output_file = "clickhouse_instances.json"
}

output "instance_count" {
  value = length(data.tencentcloud_clickhouse_instances.filtered.instance_list)
}

output "first_instance" {
  value = length(data.tencentcloud_clickhouse_instances.filtered.instance_list) > 0 ? data.tencentcloud_clickhouse_instances.filtered.instance_list[0] : null
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, String) Instance ID for exact matching, such as `cdwch-xxxxxx`.
* `instance_name` - (Optional, String) Instance name for fuzzy matching.
* `tags` - (Optional, Map) Tag filter. Multiple tags must be matched at the same time (AND logic).
* `vips` - (Optional, List) VIP address list for filtering instances.
* `is_simple` - (Optional, Bool) Whether to return simplified information. Default is `false`.
* `result_output_file` - (Optional, String) Used to save results in JSON format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - List of ClickHouse instances.
  * `access_info` - Access address, such as `10.0.0.1:9000`.
  * `can_attach_cbs_lvm` - Whether CBS LVM can be attached.
  * `can_attach_cbs` - Whether CBS can be attached.
  * `can_attach_cos` - Whether COS can be attached.
  * `ch_proxy_vip` - CHProxy VIP address.
  * `cls_log_set_id` - CLS log set ID.
  * `cls_topic_id` - CLS topic ID.
  * `common_summary` - Common node (ZooKeeper) summary information.
    * `attach_cbs_spec` - Attached CBS specification.
      * `disk_count` - Number of disks.
      * `disk_desc` - Disk description.
      * `disk_size` - Disk size in GB.
      * `disk_type` - Disk type.
    * `core` - CPU cores.
    * `disk_count` - Number of disks.
    * `disk_desc` - Disk description.
    * `disk_type` - Disk type.
    * `disk` - Disk size in GB.
    * `encrypt` - Encryption status.
    * `max_disk_size` - Maximum disk size.
    * `memory` - Memory size in GB.
    * `node_size` - Number of nodes.
    * `spec_core` - Specification CPU cores.
    * `spec_memory` - Specification memory.
    * `spec` - Specification name.
    * `sub_product_type` - Sub-product type.
  * `components` - Component list.
    * `name` - Component name.
    * `version` - Component version.
  * `cos_bucket_name` - COS bucket name.
  * `create_time` - Creation time.
  * `eip` - Elastic IP address.
  * `enable_xml_config` - Whether XML configuration is supported.
  * `expire_time` - Expiration time.
  * `flow_msg` - Workflow message.
  * `ha_zk` - ZooKeeper high availability.
  * `ha` - High availability: `true` or `false`.
  * `has_cls_topic` - Whether CLS topic is enabled.
  * `instance_id` - Instance ID, such as `cdwch-xxxx`.
  * `instance_name` - Instance name.
  * `instance_state_info` - Instance state details.
    * `flow_create_time` - Workflow creation time.
    * `flow_msg` - Workflow message.
    * `flow_name` - Workflow name.
    * `flow_progress` - Workflow progress.
    * `instance_state_desc` - Instance state description.
    * `instance_state` - Instance state.
    * `process_name` - Process name.
  * `is_elastic` - Whether it is an elastic instance.
  * `kind` - Instance type: `external`, `local`, or `yunti`.
  * `master_summary` - Master node (data node) summary information.
    * `attach_cbs_spec` - Attached CBS specification.
      * `disk_count` - Number of disks.
      * `disk_desc` - Disk description.
      * `disk_size` - Disk size in GB.
      * `disk_type` - Disk type.
    * `core` - CPU cores.
    * `disk_count` - Number of disks.
    * `disk_desc` - Disk description.
    * `disk_type` - Disk type.
    * `disk` - Disk size in GB.
    * `encrypt` - Encryption status.
    * `max_disk_size` - Maximum disk size.
    * `memory` - Memory size in GB.
    * `node_size` - Number of nodes.
    * `spec_core` - Specification CPU cores.
    * `spec_memory` - Specification memory.
    * `spec` - Specification name.
    * `sub_product_type` - Sub-product type.
  * `monitor` - Monitoring information.
  * `pay_mode` - Payment mode: `hour` (pay-as-you-go) or `prepay` (prepaid).
  * `region_desc` - Region description.
  * `region_id` - Region ID.
  * `region` - Region, such as `ap-guangzhou`.
  * `renew_flag` - Auto-renewal flag.
  * `status_desc` - Status description.
  * `status` - Instance status: `Init` (creating), `Serving` (running), `Deleted` (destroyed), `Deleting` (destroying), `Modify` (modifying).
  * `subnet_id` - Subnet ID.
  * `tags` - Tag list.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `upgrade_versions` - Upgradeable versions.
  * `version` - Instance version.
  * `vpc_id` - VPC ID.
  * `zone_desc` - Zone description.
  * `zone` - Availability zone, such as `ap-guangzhou-3`.

## Note

* This data source queries instances using the `DescribeInstancesNew` API, which has a rate limit of 20 requests per second.
* When using multiple filters, they are combined with AND logic (all conditions must be met).
* The `is_simple` parameter can be used to reduce the amount of data returned, improving query performance.
* For large numbers of instances, consider using filters to narrow down the results.
