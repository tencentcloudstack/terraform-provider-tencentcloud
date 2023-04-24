---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_cluster"
sidebar_current: "docs-tencentcloud-datasource-tsf_cluster"
description: |-
  Use this data source to query detailed information of tsf cluster
---

# tencentcloud_tsf_cluster

Use this data source to query detailed information of tsf cluster

## Example Usage

```hcl
data "tencentcloud_tsf_cluster" "cluster" {
  cluster_id_list = ["cluster-vwgj5e6y"]
  cluster_type    = "V"
  # search_word = ""
  disable_program_auth_check = true
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id_list` - (Optional, Set: [`String`]) Cluster ID list to be queried, if not filled in or passed, all content will be queried.
* `cluster_type` - (Optional, String) The type of cluster to be queried, if left blank or not passed, all content will be queried. C: container, V: virtual machine.
* `disable_program_auth_check` - (Optional, Bool) Whether to disable dataset authentication.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) Filter by keywords for Cluster Id or name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - TSF cluster pagination object. Note: This field may return null, indicating no valid value.
  * `content` - Cluster list. Note: This field may return null, indicating no valid values.
    * `cluster_cidr` - cluster CIDR. Note: This field may return null, indicating no valid value.
    * `cluster_desc` - Cluster description. Note: This field may return null, indicating no valid value.
    * `cluster_id` - Cluster ID. Note: This field may return null, indicating no valid value.
    * `cluster_limit_cpu` - Maximum CPU limit of the cluster, in cores. This field may return null, indicating that no valid value was found.
    * `cluster_limit_mem` - Cluster maximum memory limit in GB. This field may return null, indicating that no valid value was found.
    * `cluster_name` - Cluster name. Note: This field may return null, indicating no valid value.
    * `cluster_status` - cluster status. Note: This field may return null, indicating no valid value.
    * `cluster_total_cpu` - Total CPU of the cluster, unit: cores. Note: This field may return null, indicating that no valid value was found.
    * `cluster_total_mem` - Total memory of the cluster, unit: G. Note: This field may return null, indicating that no valid value is obtained.
    * `cluster_type` - Cluster type. Note: This field may return null, indicating no valid value.
    * `cluster_used_cpu` - Used CPU of the cluster, in cores. This field may return null, indicating no valid value.
    * `cluster_used_mem` - Cluster used memory in GB. This field may return null, indicating no valid value.
    * `cluster_version` - The cluster version, may return null if not applicable.
    * `create_time` - CreationTime. Note: This field may return null, indicating that no valid values can be obtained.
    * `delete_flag_reason` - Reason why the cluster cannot be deleted.  Note: This field may return null, indicating that no valid values can be obtained.
    * `delete_flag` - Deletion tag: true means it can be deleted, false means it cannot be deleted. Note: This field may return null, indicating no valid value.
    * `instance_count` - Cluster instance number. This field may return null, indicating no valid value.
    * `normal_instance_count` - Cluster normal instance number. This field may return null, indicating no valid value.
    * `operation_info` - Control information returned to the frontend. This field may return null, indicating no valid value.
      * `add_instance` - Add instance button control information, Note: This field may return null, indicating that no valid value is obtained.
        * `disabled_reason` - The reason why this button is not displayed, may return null if not applicable.
        * `enabled` - Whether the button is clickable. Note: This field may return null, indicating that no valid value is obtained.
        * `supported` - Whether the button is clickable. Note: This field may return null, indicating that no valid value was found.
      * `destroy` - Control information for destroying machine, may return null if no valid value is obtained.
        * `disabled_reason` - The reason why this button is not displayed, may return null if not applicable.
        * `enabled` - Whether the button is clickable. Note: This field may return null, indicating that no valid value is obtained.
        * `supported` - Whether the button is clickable. Note: This field may return null, indicating that no valid value was found.
      * `init` - Control information of the initialization button returned to the front end. Note: This field may return null, indicating no valid value.
        * `disabled_reason` - Reason for not displaying. Note: This field may return null, indicating no valid value.
        * `enabled` - The availability of the button (whether it is clickable) may return null indicating that the information is not available.
        * `supported` - Whether to display the button. Note: This field may return null, indicating that no valid value was found.
    * `run_instance_count` - Cluster running instance number. This field may return null, indicating no valid value.
    * `run_service_instance_count` - Number of available service instances in the cluster. Note: This field may return null, indicating no valid value.
    * `subnet_id` - Cluster subnet ID. Note: This field may return null, indicating no valid values.
    * `tsf_region_id` - region ID of TSF.  Note: This field may return null, indicating that no valid values can be obtained.
    * `tsf_region_name` - region name of TSF.  Note: This field may return null, indicating that no valid values can be obtained.
    * `tsf_zone_id` - Zone Id of TSF.  Note: This field may return null, indicating that no valid values can be obtained.
    * `tsf_zone_name` - Zone name of TSF.  Note: This field may return null, indicating that no valid values can be obtained.
    * `update_time` - last update time.  Note: This field may return null, indicating that no valid values can be obtained.
    * `vpc_id` - Private network ID of the cluster. Note: This field may return null, indicating no valid value.
  * `total_count` - Total number of items. Note: This field may return null, indicating that no valid value was found.


