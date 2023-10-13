---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_list_entities_for_policy"
sidebar_current: "docs-tencentcloud-datasource-cam_list_entities_for_policy"
description: |-
  Use this data source to query detailed information of cam list_entities_for_policy
---

# tencentcloud_cam_list_entities_for_policy

Use this data source to query detailed information of cam list_entities_for_policy

## Example Usage

```hcl
data "tencentcloud_cam_list_entities_for_policy" "list_entities_for_policy" {
  policy_id     = 1
  entity_filter = "All"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, Int) Policy Id.
* `entity_filter` - (Optional, String) Can take values of &amp;amp;#39;All&amp;amp;#39;, &amp;amp;#39;User&amp;amp;#39;, &amp;amp;#39;Group&amp;amp;#39;, and &amp;amp;#39;Role&amp;amp;#39;. &amp;amp;#39;All&amp;amp;#39; represents obtaining all entity types, &amp;amp;#39;User&amp;amp;#39; represents only obtaining sub accounts, &amp;amp;#39;Group&amp;amp;#39; represents only obtaining user groups, and &amp;amp;#39;Role&amp;amp;#39; represents only obtaining roles. The default value is&amp;amp;#39; All &amp;amp;#39;.
* `result_output_file` - (Optional, String) Used to save results.
* `rp` - (Optional, Int) Per page size, default value is 20.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Entity ListNote: This field may return null, indicating that a valid value cannot be obtained.
  * `attachment_time` - Policy association timeNote: This field may return null, indicating that a valid value cannot be obtained.
  * `id` - Entity ID.
  * `name` - Entity NameNote: This field may return null, indicating that a valid value cannot be obtained.
  * `related_type` - Association type. 1. User association; 2 User Group Association.
  * `uin` - Entity UinNote: This field may return null, indicating that a valid value cannot be obtained.


