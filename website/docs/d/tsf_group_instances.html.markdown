---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_group_instances"
sidebar_current: "docs-tencentcloud-datasource-tsf_group_instances"
description: |-
  Use this data source to query detailed information of tsf group_instances
---

# tencentcloud_tsf_group_instances

Use this data source to query detailed information of tsf group_instances

## Example Usage

```hcl
data "tencentcloud_tsf_group_instances" "group_instances" {
  group_id    = "group-yrjkln9v"
  search_word = "testing"
  order_by    = "ASC"
  order_type  = 0
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) group id.
* `order_by` - (Optional, String) order term.
* `order_type` - (Optional, Int) order type.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) search word.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Machine information of the deployment group.Note: This field may return null, which means no valid value was found.
  * `content` - List of machine instances.Note: This field may return null, which means no valid value was found.
    * `agent_version` - Agent version.Note: This field may return null, which means no valid value was found.
    * `application_id` - Application id.Note: This field may return null, which means no valid value was found.
    * `application_name` - Application name. Note: This field may return null, which means no valid value was found.
    * `application_resource_type` - application resource id.Note: This field may return null, which means no valid value was found.
    * `application_type` - Application id.Note: This field may return null, which means no valid value was found.
    * `cluster_id` - Cluster id.Note: This field may return null, which means no valid value was found.
    * `cluster_name` - Cluster name. Note: This field may return null, which means no valid value was found.
    * `cluster_type` - Cluster type.Note: This field may return null, which means no valid value was found.
    * `count_in_tsf` - Indicates whether this instance has been added to the TSF.Note: This field may return null, which means no valid value was found.
    * `group_id` - Group id.Note: This field may return null, which means no valid value was found.
    * `group_name` - Group name.Note: This field may return null, which means no valid value was found.
    * `instance_available_status` - VM availability status. For virtual machines, it indicates whether the virtual machine can be used as a resource. For containers, it indicates whether the virtual machine can be used to deploy pods.Note: This field may return null, which means no valid value was found.
    * `instance_charge_type` - machine instance charge type.Note: This field may return null, which means no valid value was found.
    * `instance_created_time` - Creation time of the machine instance in CVM.Note: This field may return null, which means no valid value was found.
    * `instance_desc` - Description.Note: This field may return null, which means no valid value was found.
    * `instance_expired_time` - Expire time of the machine instance in CVM.Note: This field may return null, which means no valid value was found.
    * `instance_id` - Machine instance ID.Note: This field may return null, which means no valid value was found.
    * `instance_import_mode` - InstanceImportMode import mode.Note: This field may return null, which means no valid value was found.
    * `instance_limit_cpu` - Limit CPU information of the machine instance.Note: This field may return null, which means no valid value was found.
    * `instance_limit_mem` - Limit memory information of the machine instance.Note: This field may return null, which means no valid value was found.
    * `instance_name` - Machine name.Note: This field may return null, which means no valid value was found.
    * `instance_pkg_version` - instance pkg version.Note: This field may return null, which means no valid value was found.
    * `instance_status` - VM status. For virtual machines, it indicates the status of the virtual machine. For containers, it indicates the status of the virtual machine where the pod is located.Note: This field may return null, which means no valid value was found.
    * `instance_total_cpu` - Total CPU information of the machine instance.Note: This field may return null, which means no valid value was found.
    * `instance_total_mem` - Total memory information of the machine instance.Note: This field may return null, which means no valid value was found.
    * `instance_used_cpu` - CPU information used by the machine instance.Note: This field may return null, which means no valid value was found.
    * `instance_used_mem` - Memory information used by the machine instance.Note: This field may return null, which means no valid value was found.
    * `instance_zone_id` - Instance zone id.Note: This field may return null, which means no valid value was found.
    * `lan_ip` - Private IP address.Note: This field may return null, which means no valid value was found.
    * `namespace_id` - Namespace id.Note: This field may return null, which means no valid value was found.
    * `namespace_name` - Namespace name.Note: This field may return null, which means no valid value was found.
    * `node_instance_id` - Container host instance ID.Note: This field may return null, which means no valid value was found.
    * `operation_state` - Execution status of the instance.Note: This field may return null, which means no valid value was found.
    * `reason` - Health checking reason.Note: This field may return null, which means no valid value was found.
    * `restrict_state` - Business status of the machine instance.Note: This field may return null, which means no valid value was found.
    * `service_instance_status` - Status of service instances under the service. For virtual machines, it indicates whether the application is available and the agent status. For containers, it indicates the status of the pod.Note: This field may return null, which means no valid value was found.
    * `service_sidecar_status` - Sidecar status.Note: This field may return null, which means no valid value was found.
    * `update_time` - Update time.Note: This field may return null, which means no valid value was found.
    * `wan_ip` - Public IP address.Note: This field may return null, which means no valid value was found.
  * `total_count` - Total number of machine instances.Note: This field may return null, which means no valid value was found.


