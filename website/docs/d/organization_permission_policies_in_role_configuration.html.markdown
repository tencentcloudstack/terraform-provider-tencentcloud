---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_permission_policies_in_role_configuration"
sidebar_current: "docs-tencentcloud-datasource-organization_permission_policies_in_role_configuration"
description: |-
  Use this data source to query detailed information of Organization permission policies in role configuration
---

# tencentcloud_organization_permission_policies_in_role_configuration

Use this data source to query detailed information of Organization permission policies in role configuration

## Example Usage

### Query all permission policies in a role configuration

```hcl
data "tencentcloud_organization_permission_policies_in_role_configuration" "example" {
  zone_id               = "z-xxxxxx"
  role_configuration_id = "rc-xxxxxx"
}
```

### Query permission policies filtered by policy type

```hcl
data "tencentcloud_organization_permission_policies_in_role_configuration" "example" {
  zone_id               = "z-xxxxxx"
  role_configuration_id = "rc-xxxxxx"
  role_policy_type      = "System"
}
```

### Query permission policies filtered by policy name

```hcl
data "tencentcloud_organization_permission_policies_in_role_configuration" "example" {
  zone_id               = "z-xxxxxx"
  role_configuration_id = "rc-xxxxxx"
  filter                = "AdministratorAccess"
}
```

## Argument Reference

The following arguments are supported:

* `role_configuration_id` - (Required, String) Role configuration ID.
* `zone_id` - (Required, String) Space ID.
* `filter` - (Optional, String) Search by policy name.
* `result_output_file` - (Optional, String) Used to save results.
* `role_policy_type` - (Optional, String) Permission policy type. Valid values: `System`: System policy, reuses CAM system policies. `Custom`: Custom policy, written according to CAM permission policy syntax and structure.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `role_policies` - Permission policy list.
  * `add_time` - Time when the permission policy was added to the role configuration.
  * `role_policy_document` - Custom policy content. Only returned for custom policies.
  * `role_policy_id` - Policy ID.
  * `role_policy_name` - Permission policy name.
  * `role_policy_type` - Permission policy type.
* `total_counts` - Total number of permission policies.


