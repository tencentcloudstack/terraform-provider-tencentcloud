---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_role_configuration"
sidebar_current: "docs-tencentcloud-resource-identity_center_role_configuration"
description: |-
  Provides a resource to create a organization identity_center_role_configuration
---

# tencentcloud_identity_center_role_configuration

Provides a resource to create a organization identity_center_role_configuration

## Example Usage

```hcl
resource "tencentcloud_identity_center_role_configuration" "identity_center_role_configuration" {
  zone_id                 = "z-xxxxxx"
  role_configuration_name = "tf-test"
  description             = "test"
}
```

## Argument Reference

The following arguments are supported:

* `role_configuration_name` - (Required, String) Access configuration name, which contains up to 128 characters, including English letters, digits, and hyphens (-).
* `zone_id` - (Required, String) Space ID.
* `description` - (Optional, String) Access configuration description, which contains up to 1024 characters.
* `relay_state` - (Optional, String) Initial access page. It indicates the initial access page URL when CIC users use the access configuration to access the target account of the Tencent Cloud Organization. This page must be the Tencent Cloud console page. The default is null, which indicates navigating to the home page of the Tencent Cloud console.
* `session_duration` - (Optional, Int) Session duration. It indicates the maximum session duration when CIC users use the access configuration to access the target account of the Tencent Cloud Organization. Unit: seconds. Value range: 900-43,200 (15 minutes to 12 hours). Default value: 3600 (1 hour).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `role_configuration_id` - Role configuration id.
* `update_time` - Update time.


## Import

organization identity_center_role_configuration can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_role_configuration.identity_center_role_configuration ${zoneId}#${roleConfigurationId}
```

