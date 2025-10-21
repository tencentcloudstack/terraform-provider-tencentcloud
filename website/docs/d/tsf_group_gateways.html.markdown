---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_group_gateways"
sidebar_current: "docs-tencentcloud-datasource-tsf_group_gateways"
description: |-
  Use this data source to query detailed information of tsf group_gateways
---

# tencentcloud_tsf_group_gateways

Use this data source to query detailed information of tsf group_gateways

## Example Usage

```hcl
data "tencentcloud_tsf_group_gateways" "group_gateways" {
  gateway_deploy_group_id = "group-aeoej4qy"
  search_word             = "test"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_deploy_group_id` - (Required, String) gateway group Id.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) search word.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - api group information.
  * `content` - api group Info.
    * `acl_mode` - ACL type for accessing the group.Note: This field may return null, which means no valid value was found.
    * `api_count` - Number of APIs.Note: This field may return null, which means no valid value was found.
    * `auth_type` - Authentication type. secret: key authentication; none: no authentication.Note: This field may return null, which means no valid value was found.
    * `binded_gateway_deploy_groups` - Gateway deployment group bound to the API group.Note: This field may return null, which means no valid value was found.
      * `application_id` - application ID.Note: This field may return null, which means no valid value was found.
      * `application_name` - application name.Note: This field may return null, which means no valid value was found.
      * `application_type` - Application category: V: virtual machine application, C: container application.Note: This field may return null, which means no valid value was found.
      * `cluster_type` - Cluster type, with possible values: C: container, V: virtual machine.Note: This field may return null, which means no valid value was found.
      * `deploy_group_id` - Gateway deployment group ID.Note: This field may return null, which means no valid value was found.
      * `deploy_group_name` - Gateway deployment group name.Note: This field may return null, which means no valid value was found.
      * `group_status` - Application status of the deployment group, with possible values: Running, Waiting, Paused, Updating, RollingBack, Abnormal, Unknown.Note: This field may return null, which means no valid value was found.
    * `created_time` - Group creation time, such as: 2019-06-20 15:51:28.Note: This field may return null, which means no valid value was found.
    * `description` - Description.Note: This field may return null, which means no valid value was found.
    * `gateway_instance_id` - Gateway instance ID.Note: This field may return null, which means no valid value was found.
    * `gateway_instance_type` - Gateway instance type.Note: This field may return null, which means no valid value was found.
    * `group_context` - api group context.Note: This field may return null, which means no valid value was found.
    * `group_id` - api group id.Note: This field may return null, which means no valid value was found.
    * `group_name` - api group name.Note: This field may return null, which means no valid value was found.
    * `group_type` - Group type. ms: microservice group; external: external API group.This field may return null, which means no valid value was found.
    * `namespace_name_key_position` - Namespace parameter location, path, header, or query. The default is path.Note: This field may return null, which means no valid value was found.
    * `namespace_name_key` - Namespace parameter key.Note: This field may return null, which means no valid value was found.
    * `service_name_key_position` - Microservice name parameter location, path, header, or query. The default is path.Note: This field may return null, which means no valid value was found.
    * `service_name_key` - Microservice name parameter key.Note: This field may return null, which means no valid value was found.
    * `status` - Release status. drafted: not released. released: released.Note: This field may return null, which means no valid value was found.
    * `updated_time` - Group update time, such as: 2019-06-20 15:51:28.Note: This field may return null, which means no valid value was found.
  * `total_count` - total count.


