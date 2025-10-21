---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_auth_node"
sidebar_current: "docs-tencentcloud-datasource-organization_org_auth_node"
description: |-
  Use this data source to query detailed information of organization org_auth_node
---

# tencentcloud_organization_org_auth_node

Use this data source to query detailed information of organization org_auth_node

## Example Usage

```hcl
data "tencentcloud_organization_org_auth_node" "org_auth_node" {
}
```

## Argument Reference

The following arguments are supported:

* `auth_name` - (Optional, String) Verified company name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Organization auth node list.
  * `auth_name` - Verified company name.
  * `manager` - Organization auth manager.
    * `member_name` - Member name.
    * `member_uin` - Member uin.
  * `relation_id` - Relationship Id.


