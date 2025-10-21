---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_api_group"
sidebar_current: "docs-tencentcloud-resource-tsf_api_group"
description: |-
  Provides a resource to create a tsf api_group
---

# tencentcloud_tsf_api_group

Provides a resource to create a tsf api_group

## Example Usage

```hcl
resource "tencentcloud_tsf_api_group" "api_group" {
  group_name          = "terraform_test_group"
  group_context       = "/terraform-test"
  auth_type           = "none"
  description         = "terraform-test"
  group_type          = "ms"
  gateway_instance_id = "gw-ins-i6mjpgm8"
  # namespace_name_key = "path"
  # service_name_key = "path"
  namespace_name_key_position = "path"
  service_name_key_position   = "path"
}
```

## Argument Reference

The following arguments are supported:

* `group_context` - (Required, String) grouping context.
* `group_name` - (Required, String) group name, cannot contain Chinese.
* `auth_type` - (Optional, String) authentication type. secret: key authentication; none: no authentication.
* `description` - (Optional, String) remarks.
* `gateway_instance_id` - (Optional, String) gateway entity ID.
* `group_type` - (Optional, String) grouping type, default ms. ms: microservice grouping; external: external Api grouping.
* `namespace_name_key_position` - (Optional, String) namespace parameter position, path, header or query, the default is path.
* `namespace_name_key` - (Optional, String) namespace parameter key value.
* `service_name_key_position` - (Optional, String) microservice name parameter position, path, header or query, the default is path.
* `service_name_key` - (Optional, String) microservice name parameter key value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `acl_mode` - Access group ACL type.
* `api_count` - number of APIs.
* `binded_gateway_deploy_groups` - api group bound gateway deployment group.
  * `application_id` - application ID.
  * `application_name` - Application Name.
  * `application_type` - Application classification: V: virtual machine application, C: container application.
  * `cluster_type` - Cluster type, C: container, V: virtual machine.
  * `deploy_group_id` - Gateway deployment group ID.
  * `deploy_group_name` - Gateway deployment group name.
  * `group_status` - Deployment group application status, values: Running, Waiting, Paused, Updating, RollingBack, Abnormal, Unknown.
* `created_time` - Group creation time such as: 2019-06-20 15:51:28.
* `gateway_instance_type` - Type of gateway instance.
* `group_id` - Api Group Id.
* `status` - Release status, drafted: Not published. released: released.
* `updated_time` - Group update time such as: 2019-06-20 15:51:28.


## Import

tsf api_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_api_group.api_group api_group_id
```

