---
subcategory: "EMR"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_nodes"
sidebar_current: "docs-tencentcloud-datasource-emr_nodes"
description: |-
  Provides an available EMR for the user.
---

# tencentcloud_emr_nodes

Provides an available EMR for the user.

The EMR data source obtain the hardware node information by using the emr cluster ID.

## Example Usage

```hcl
data "tencentcloud_emr_nodes" "my_emr_nodes" {
  node_flag   = "master"
  instance_id = "emr-rnzqrleq"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Cluster instance ID, the instance ID is as follows: emr-xxxxxxxx.
* `node_flag` - (Required) Node ID, the value is:
				- all: Means to get all type nodes, except cdb information.
				- master: Indicates that the master node information is obtained.
				- core: Indicates that the core node information is obtained.
				- task: indicates obtaining task node information.
				- common: means to get common node information.
				- router: Indicates obtaining router node information.
				- db: Indicates that the cdb information for the normal state is obtained.
				- recyle: Indicates that the node information in the Recycle Bin isolation, including the cdb information, is obtained.
				- renew: Indicates that all node information to be renewed, including cddb information, is obtained, and the auto-renewal node will not be returned.
				
				Note: Only the above values are now supported, entering other values will cause an error.
* `hardware_resource_type` - (Optional) Resource type: Support all/host/pod, default is all.
* `limit` - (Optional) The number returned per page, the default value is 100, and the maximum value is 100.
* `offset` - (Optional) Page number, with a default value of 0, represents the first page.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nodes` - List of node details.
  * `app_id` - User APPID.
  * `apply_time` - Application time.
  * `auto_flag` - Whether it is an autoscaling node, 0 is a normal node, and 1 is an autoscaling node.
  * `cdb_ip` - Database IP.
  * `cdb_node_info` - Database information.
    * `apply_time` - Application time.
    * `expire_flag` - Expired id.
    * `expire_time` - Expiration time.
    * `instance_name` - DB instance.
    * `ip` - Database IP.
    * `is_auto_renew` - Renewal identity.
    * `mem_size` - Database memory specifications.
    * `pay_type` - The type of payment.
    * `port` - Database port.
    * `region_id` - Region id.
    * `serial_no` - Database string.
    * `service` - The service identity.
    * `status` - Database status.
    * `volume` - Database disk specifications.
    * `zone_id` - Zone Id.
  * `cdb_port` - Database port.
  * `charge_type` - The type of payment.
  * `cpu_num` - Number of node cores.
  * `destroyable` - Whether this node is destroyable, 1 can be destroyed, 0 is not destroyable.
  * `device_class` - Device identity.
  * `disk_size` - Hard disk size.
  * `dynamic_pod_spec` - Floating specification value json string.
  * `emr_resource_id` - Node resource ID.
  * `expire_time` - Expiration time.
  * `flag` - Node type. 0: common node; 1: master node; 2: core node; 3: task node.
  * `free_time` - Release time.
  * `hardware_resource_type` - Resource type, host/pod.
  * `hw_disk_size_desc` - Hard disk capacity description.
  * `hw_disk_size` - Hard disk capacity.
  * `hw_mem_size_desc` - Memory capacity description.
  * `hw_mem_size` - Memory capacity.
  * `ip` - Intranet IP.
  * `is_auto_renew` - Renewal logo.
  * `is_dynamic_spec` - Floating specifications, 1 yes, 0 no.
  * `mc_multi_disks` - Multi-cloud disk.
    * `count` - The number of cloud disks of this type.
    * `type` - Disk type.
    * `volume` - The size of the cloud disk.
  * `mem_desc` - Node memory description.
  * `mem_size` - Node memory.
  * `mutable` - Supports variations.
  * `name_tag` - Node description.
  * `order_no` - Machine instance ID.
  * `region_id` - The node is located in the region.
  * `root_size` - The size of the system disk.
  * `serial_no` - Serial number.
  * `services` - Node deployment service.
  * `spec` - Node specifications.
  * `storage_type` - Disk type.
  * `support_modify_pay_mode` - Whether to support change billing type 1 Yes and 0 No.
  * `tags` - The label of the node binding.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `wan_ip` - The master node is bound to the Internet IP address.
  * `zone_id` - Zone where the node is located.


