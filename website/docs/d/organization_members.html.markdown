---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_members"
sidebar_current: "docs-tencentcloud-datasource-organization_members"
description: |-
  Use this data source to query detailed information of organization members
---

# tencentcloud_organization_members

Use this data source to query detailed information of organization members

## Example Usage

```hcl
data "tencentcloud_organization_members" "members" {}
```

## Argument Reference

The following arguments are supported:

* `auth_name` - (Optional, String) Entity name.
* `lang` - (Optional, String) Valid values: `en` (Tencent Cloud International); `zh` (Tencent Cloud).
* `product` - (Optional, String) Abbreviation of the trusted service, which is required during querying the trusted service admin.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Search by member name or ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Member list.
  * `bind_status` - Security information binding status. Valid values: `Unbound`, `Valid`, `Success`, `Failed`.Note: This field may return null, indicating that no valid values can be obtained.
  * `create_time` - Creation timeNote: This field may return null, indicating that no valid values can be obtained.
  * `is_allow_quit` - Whether the member is allowed to leave. Valid values: `Allow`, `Denied`.Note: This field may return null, indicating that no valid values can be obtained.
  * `member_type` - Member type. Valid values: `Invite` (invited); `Create` (created).Note: This field may return null, indicating that no valid values can be obtained.
  * `member_uin` - Member UINNote: This field may return null, indicating that no valid values can be obtained.
  * `name` - Member nameNote: This field may return null, indicating that no valid values can be obtained.
  * `node_id` - Node IDNote: This field may return null, indicating that no valid values can be obtained.
  * `node_name` - Node nameNote: This field may return null, indicating that no valid values can be obtained.
  * `org_identity` - Management identityNote: This field may return null, indicating that no valid values can be obtained.
    * `identity_alias_name` - Identity name.Note: This field may return null, indicating that no valid values can be obtained.
    * `identity_id` - Identity ID.Note: This field may return null, indicating that no valid values can be obtained.
  * `org_permission` - Relationship policy permissionNote: This field may return null, indicating that no valid values can be obtained.
    * `id` - Permission ID.
    * `name` - Permission name.
  * `org_policy_name` - Relationship policy nameNote: This field may return null, indicating that no valid values can be obtained.
  * `org_policy_type` - Relationship policy typeNote: This field may return null, indicating that no valid values can be obtained.
  * `pay_name` - Payer nameNote: This field may return null, indicating that no valid values can be obtained.
  * `pay_uin` - Payer UINNote: This field may return null, indicating that no valid values can be obtained.
  * `permission_status` - Member permission status. Valid values: `Confirmed`, `UnConfirmed`.Note: This field may return null, indicating that no valid values can be obtained.
  * `remark` - RemarksNote: This field may return null, indicating that no valid values can be obtained.
  * `update_time` - Update timeNote: This field may return null, indicating that no valid values can be obtained.


