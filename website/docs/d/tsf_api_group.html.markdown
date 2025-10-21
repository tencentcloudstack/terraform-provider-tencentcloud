---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_api_group"
sidebar_current: "docs-tencentcloud-datasource-tsf_api_group"
description: |-
  Use this data source to query detailed information of tsf api_group
---

# tencentcloud_tsf_api_group

Use this data source to query detailed information of tsf api_group

## Example Usage

```hcl
data "tencentcloud_tsf_api_group" "api_group" {
  search_word         = "xxx01"
  group_type          = "ms"
  auth_type           = "none"
  status              = "released"
  order_by            = "created_time"
  order_type          = 0
  gateway_instance_id = "gw-ins-lvdypq5k"
}
```

## Argument Reference

The following arguments are supported:

* `auth_type` - (Optional, String) Authentication type. secret: Secret key authentication; none: No authentication.
* `gateway_instance_id` - (Optional, String) Gateway Instance Id.
* `group_type` - (Optional, String) Group type. ms: Microservice group; external: External API group.
* `order_by` - (Optional, String) Sorting field: created_time or group_context.
* `order_type` - (Optional, Int) Sorting type: 0 (ASC) or 1 (DESC).
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) search word.
* `status` - (Optional, String) Publishing status. drafted: Not published. released: Published.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Pagination structure.Note: This field may return null, indicating that no valid values can be obtained.
  * `content` - Api group info.
    * `acl_mode` - Number of APIs.Note: This field may return null, indicating that no valid values can be obtained.
    * `api_count` - api count.
    * `auth_type` - Authentication type. secret: key authentication; none: no authentication.Note: This field may return null, indicating that no valid values can be obtained.
    * `binded_gateway_deploy_groups` - The gateway group bind with the api group list.
      * `application_id` - Application ID.Note: This field may return null, indicating that no valid values can be obtained.
      * `application_name` - Application Name.Note: This field may return null, indicating that no valid values can be obtained.
      * `application_type` - Application Name.Note: This field may return null, indicating that no valid values can be obtained.
      * `cluster_type` - Cluster type, C: container, V: virtual machine.Note: This field may return null, indicating that no valid values can be obtained.
      * `deploy_group_id` - Gateway deployment group bound to the API group.Note: This field may return null, indicating that no valid values can be obtained.
      * `deploy_group_name` - Deploy group name.Note: This field may return null, indicating that no valid values can be obtained.
      * `group_status` - Application category: V: virtual machine application, C: container application. Note: This field may return null, indicating that no valid values can be obtained.
    * `created_time` - Group creation time.Note: This field may return null, indicating that no valid values can be obtained.
    * `description` - Description.Note: This field may return null, indicating that no valid values can be obtained.
    * `gateway_instance_id` - Gateway Instance Id.Note: This field may return null, indicating that no valid values can be obtained.
    * `gateway_instance_type` - Gateway Instance Type.Note: This field may return null, indicating that no valid values can be obtained.
    * `group_context` - Api Group Context.Note: This field may return null, indicating that no valid values can be obtained.
    * `group_id` - Api Group Id.Note: This field may return null, indicating that no valid values can be obtained.
    * `group_name` - Api Group Name.Note: This field may return null, indicating that no valid values can be obtained.
    * `group_type` - Group type.Note: This field may return null, indicating that no valid values can be obtained.
    * `namespace_name_key_position` - Namespace parameter location, path, header, or query, default is path. Note: This field may return null, indicating that no valid values can be obtained.
    * `namespace_name_key` - Namespace name key.Note: This field may return null, indicating that no valid values can be obtained.
    * `service_name_key_position` - Microservice name parameter location, path, header, or query, default is path.Note: This field may return null, indicating that no valid values can be obtained.
    * `service_name_key` - Key value of microservice name parameter.Note: This field may return null, indicating that no valid values can be obtained.
    * `status` - Release status. drafted: not released. released: released.Note: This field may return null, indicating that no valid values can be obtained.
    * `updated_time` - Group creation time, such as: 2019-06-20 15:51:28.Note: This field may return null, indicating that no valid values can be obtained.
  * `total_count` - record count.


