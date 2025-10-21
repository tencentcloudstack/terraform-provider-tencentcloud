---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_groups"
sidebar_current: "docs-tencentcloud-datasource-tsf_groups"
description: |-
  Use this data source to query detailed information of tsf groups
---

# tencentcloud_tsf_groups

Use this data source to query detailed information of tsf groups

## Example Usage

```hcl
data "tencentcloud_tsf_groups" "groups" {
  search_word              = "keep"
  application_id           = "application-a24x29xv"
  order_by                 = "createTime"
  order_type               = 0
  namespace_id             = "namespace-aemrg36v"
  cluster_id               = "cluster-vwgj5e6y"
  group_resource_type_list = ["DEF"]
  status                   = "Running"
  group_id_list            = ["group-yrjkln9v"]
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Optional, String) applicationId.
* `cluster_id` - (Optional, String) clusterId.
* `group_id_list` - (Optional, Set: [`String`]) group Id list.
* `group_resource_type_list` - (Optional, Set: [`String`]) Group resourceType list.
* `namespace_id` - (Optional, String) namespace Id.
* `order_by` - (Optional, String) sort term.
* `order_type` - (Optional, Int) order type, 0 desc, 1 asc.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) searchWord, support groupName.
* `status` - (Optional, String) group status filter, `Running`: running, `Unknown`: unknown.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Pagination information of the virtual machine deployment group.Note: This field may return null, indicating that no valid value was found.
  * `content` - Virtual machine deployment group list. Note: This field may return null, indicating that no valid value was found.
    * `alias` - Group alias. Note: This field may return null, indicating that no valid value was found.
    * `application_id` - Application ID. Note: This field may return null, indicating that no valid value was found.
    * `application_name` - Application name. Note: This field may return null, indicating that no valid value was found.
    * `application_type` - Application type. Note: This field may return null, indicating that no valid value was found.
    * `cluster_id` - Cluster ID. Note: This field may return null, indicating that no valid value was found.
    * `cluster_name` - Cluster name. Note: This field may return null, indicating that no valid value was found.
    * `create_time` - Create Time. Note: This field may return null, indicating that no valid value was found.
    * `deploy_desc` - Group description. Note: This field may return null, indicating that no valid value was found.
    * `group_desc` - Group description. Note: This field may return null, indicating that no valid value was found.
    * `group_id` - Group ID. Note: This field may return null, indicating that no valid value was found.
    * `group_name` - Group ID. Note: This field may return null, indicating that no valid value was found.
    * `group_resource_type` - Group resource type. Note: This field may return null, indicating that no valid value was found.
    * `microservice_type` - Microservice type. Note: This field may return null, indicating that no valid value was found.
    * `namespace_id` - Namespace ID. Note: This field may return null, indicating that no valid value was found.
    * `namespace_name` - Namespace name. Note: This field may return null, indicating that no valid value was found.
    * `startup_parameters` - Group start up Parameters. Note: This field may return null, indicating that no valid value was found.
    * `update_time` - Group update time. Note: This field may return null, indicating that no valid value was found.
    * `updated_time` - Update time. Note: This field may return null, indicating that no valid value was found.
  * `total_count` - Total count virtual machine deployment group. Note: This field may return null, indicating that no valid value was found.


