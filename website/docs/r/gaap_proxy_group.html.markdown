---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_proxy_group"
sidebar_current: "docs-tencentcloud-resource-gaap_proxy_group"
description: |-
  Provides a resource to create a gaap proxy group
---

# tencentcloud_gaap_proxy_group

Provides a resource to create a gaap proxy group

## Example Usage

```hcl
resource "tencentcloud_gaap_proxy_group" "proxy_group" {
  project_id         = 0
  group_name         = "tf-test-update"
  real_server_region = "Beijing"
  ip_address_version = "IPv4"
  package_type       = "Thunder"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required, String) Channel group alias.
* `project_id` - (Required, Int) ID of the project to which the proxy group belongs.
* `real_server_region` - (Required, String) real server region, refer to the interface DescribeDestRegions to return the RegionId in the parameter RegionDetail.
* `ip_address_version` - (Optional, String) IP version, can be taken as IPv4 or IPv6 with a default value of IPv4.
* `package_type` - (Optional, String) Package type of channel group. Available values: Thunder and Accelerator. Default is Thunder.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

gaap proxy_group can be imported using the id, e.g.

```
terraform import tencentcloud_gaap_proxy_group.proxy_group proxy_group_id
```

