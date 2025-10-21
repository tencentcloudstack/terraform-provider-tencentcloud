---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_bind_api_group"
sidebar_current: "docs-tencentcloud-resource-tsf_bind_api_group"
description: |-
  Provides a resource to create a tsf bind_api_group
---

# tencentcloud_tsf_bind_api_group

Provides a resource to create a tsf bind_api_group

## Example Usage

```hcl
resource "tencentcloud_tsf_bind_api_group" "bind_api_group" {
  gateway_deploy_group_id = "group-vzd97zpy"
  group_id                = "grp-qp0rj3zi"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_deploy_group_id` - (Required, String, ForceNew) gateway group id.
* `group_id` - (Required, String, ForceNew) group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tsf bind_api_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_bind_api_group.bind_api_group bind_api_group_id
```

