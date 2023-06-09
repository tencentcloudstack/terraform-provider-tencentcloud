---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_gateway_all_group_apis"
sidebar_current: "docs-tencentcloud-datasource-tsf_gateway_all_group_apis"
description: |-
  Use this data source to query detailed information of tsf gateway_all_group_apis
---

# tencentcloud_tsf_gateway_all_group_apis

Use this data source to query detailed information of tsf gateway_all_group_apis

## Example Usage

```hcl
data "tencentcloud_tsf_gateway_all_group_apis" "gateway_all_group_apis" {
  gateway_deploy_group_id = "group-aeoej4qy"
  search_word             = "user"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_deploy_group_id` - (Required, String) gateway group Id.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) Search keyword, supports api group name or API path.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Gateway group and API list information.
  * `gateway_deploy_group_id` - Gateway deployment group ID.Note: This field may return null, which means no valid value was found.
  * `gateway_deploy_group_name` - Gateway deployment group name.Note: This field may return null, which means no valid value was found.
  * `group_num` - Gateway deployment api group number.Note: This field may return null, which means no valid value was found.
  * `groups` - Gateway deployment  api group list.Note: This field may return null, which means no valid value was found.
    * `gateway_instance_id` - gateway instance id.Note: This field may return null, which means no valid value was found.
    * `gateway_instance_type` - Type of the gateway instance.Note: This field may return null, which means no valid value was found.
    * `group_api_count` - Number of APIs under the group. Note: This field may return null, which means no valid value was found.
    * `group_apis` - List of APIs under the group.Note: This field may return null, which means no valid value was found.
      * `api_id` - API ID.
      * `method` - API method.
      * `microservice_name` - API service name.
      * `namespace_name` - namespace name.
      * `path` - API path.
    * `group_id` - api group id.Note: This field may return null, which means no valid value was found.
    * `group_name` - api group name.Note: This field may return null, which means no valid value was found.


