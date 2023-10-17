---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_work_spaces"
sidebar_current: "docs-tencentcloud-datasource-oceanus_work_spaces"
description: |-
  Use this data source to query detailed information of oceanus work_spaces
---

# tencentcloud_oceanus_work_spaces

Use this data source to query detailed information of oceanus work_spaces

## Example Usage

```hcl
data "tencentcloud_oceanus_work_spaces" "example" {
  order_type = 1
  filters {
    name   = "WorkSpaceName"
    values = ["tf_example"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter rules.
* `order_type` - (Optional, Int) 1:sort by creation time in descending order (default); 2:sort by creation time in ascending order; 3:sort by status in descending order; 4:sort by status in ascending order; default is 0.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Filter values for the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `work_space_set_item` - List of workspace detailsNote: This field may return null, indicating that no valid values can be obtained.
  * `app_id` - User APPID.
  * `cluster_group_set_item` - Workspace cluster information.
    * `app_id` - Account APPID.
    * `cluster_id` - SerialId of the clusterGroup.
    * `create_time` - Cluster creation time.
    * `creator_uin` - Creator account UIN.
    * `cu_mem` - CU memory specification.
    * `cu_num` - CU quantity.
    * `free_cu_num` - Free CU.
    * `free_cu` - Free CU under fine-grained resources.
    * `name` - Cluster name.
    * `net_environment_type` - Network.
    * `owner_uin` - Main account UIN.
    * `pay_mode` - Payment mode.
    * `region` - Region.
    * `remark` - Description.
    * `running_cu` - Running CU.
    * `status_desc` - Status description.
    * `status` - Cluster status, 1:uninitialized, 3:initializing, 2:running.
    * `update_time` - Last operation time on the cluster.
    * `zone` - Zone.
  * `create_time` - Creation time.
  * `creator_uin` - Creator UIN.
  * `description` - Workspace description.
  * `jobs_count` - Note: This field may return null, indicating that no valid values can be obtained.
  * `owner_uin` - Main account UIN.
  * `region` - Region.
  * `role_auth_count` - Workspace member count.
  * `role_auth` - Workspace role information.
    * `app_id` - User AppID.
    * `auth_sub_account_uin` - Bound authorized UIN.
    * `create_time` - Creation time.
    * `creator_uin` - Creator UIN.
    * `id` - IDNote: This field may return null, indicating that no valid values can be obtained.
    * `owner_uin` - Main account UIN.
    * `permission` - Corresponding to the ID in the role table.
    * `role_name` - Permission nameNote: This field may return null, indicating that no valid values can be obtained.
    * `status` - 2:enabled, 1:disabled.
    * `update_time` - Last operation time.
    * `work_space_id` - Workspace IDNote: This field may return null, indicating that no valid values can be obtained.
    * `work_space_serial_id` - Workspace SerialId.
  * `serial_id` - Workspace SerialId.
  * `status` - 1:uninitialized; 2:available; -1:deleted.
  * `update_time` - Update time.
  * `work_space_id` - Workspace SerialId.
  * `work_space_name` - Workspace name.


