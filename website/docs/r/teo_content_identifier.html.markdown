---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_content_identifier"
sidebar_current: "docs-tencentcloud-resource-teo_content_identifier"
description: |-
  Provides a resource to create a TEO content identifier
---

# tencentcloud_teo_content_identifier

Provides a resource to create a TEO content identifier

## Example Usage

```hcl
resource "tencentcloud_teo_content_identifier" "example" {
  plan_id     = "edgeone-6bzvsgjkfa9g"
  description = "example"
  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required, String) Description of the content identifier, length limit of up to 20 characters.
* `plan_id` - (Required, String) Target plan id to be bound, available only for the enterprise edition. <li>if there is already a plan under your account, go to [plan management](https://console.cloud.tencent.com/edgeone/package) to get the plan id and directly bind the content identifier to the plan;</li><li>if you do not have a plan to bind, please purchase an enterprise edition plan first.</li>.
* `tags` - (Optional, List) Tags of the content identifier. this parameter is used for authority control. to create tags, go to the [tag console](https://console.cloud.tencent.com/tag/taglist).

The `tags` object supports the following:

* `tag_key` - (Required, String) The tag key.
Note: This field may return null, indicating that no valid values can be obtained.
* `tag_value` - (Required, String) The tag value.
Note: This field may return null, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `content_id` - Content identifier ID.
* `created_on` - Creation time, which is in Coordinated Universal Time (UTC) and follows the ISO 8601 date and time format..
* `modified_on` - The time of the latest update, in Coordinated Universal Time (UTC), following the ISO 8601 date and time format..
* `status` - Content identifier status. Valid values: `active` (effective), `deleted` (deleted).

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateContentIdentifier](https://cloud.tencent.com/document/api/1552/114249) - This API is used to create content identifiers, where you can set descriptions, tags, and other information. It is also necessary to bind an enterprise edition package for billing data statistics. A content identifier can only bind one billing package, while a billing package can bind multiple content identifiers. This feature is only available to the allowlist.
* Read: [DescribeContentIdentifiers](https://cloud.tencent.com/document/api/1552/114247) - Batch query content identifiers, which can be filtered by ID, description, status, or Tag. Deleted content identifiers queried by status are retained for only three months. This feature is only open to the allowlist.
* Update: [ModifyContentIdentifier](https://cloud.tencent.com/document/api/1552/114246) - Modify content identifier, only description modification is supported. This feature is only open to the allowlist.
* Delete: [DeleteContentIdentifier](https://cloud.tencent.com/document/api/1552/114248) - Delete the specified content identifier. This feature is only available to the allowlist.

## Import

TEO content identifier can be imported using the id, e.g.

```
terraform import tencentcloud_teo_content_identifier.example eocontent-3dy8iyfq8dba
```

