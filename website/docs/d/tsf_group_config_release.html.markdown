---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_group_config_release"
sidebar_current: "docs-tencentcloud-datasource-tsf_group_config_release"
description: |-
  Use this data source to query detailed information of tsf group_config_release
---

# tencentcloud_tsf_group_config_release

Use this data source to query detailed information of tsf group_config_release

## Example Usage

```hcl
data "tencentcloud_tsf_group_config_release" "group_config_release" {
  group_id = "group-yrjkln9v"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) groupId.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Information related to the deployment group release.Note: This field may return null, which means no valid value was found.
  * `config_release_list` - Configuration item release list.Note: This field may return null, which means no valid value was found.
    * `application_id` - Configuration item release config ID.Note: This field may return null, which means no valid value was found.
    * `cluster_id` - Configuration item release cluster ID.Note: This field may return null, which means no valid value was found.
    * `cluster_name` - Configuration item release cluster name.Note: This field may return null, which means no valid value was found.
    * `config_id` - Configuration item release config ID.Note: This field may return null, which means no valid value was found.
    * `config_name` - Configuration item release config name.Note: This field may return null, which means no valid value was found.
    * `config_release_id` - Configuration item release ID.Note: This field may return null, which means no valid value was found.
    * `config_version` - Configuration item release config version.Note: This field may return null, which means no valid value was found.
    * `group_id` - Configuration item release config group ID.Note: This field may return null, which means no valid value was found.
    * `group_name` - Configuration item release config group name.Note: This field may return null, which means no valid value was found.
    * `namespace_id` - Configuration item release namespace ID.Note: This field may return null, which means no valid value was found.
    * `namespace_name` - Configuration item release namespace name.Note: This field may return null, which means no valid value was found.
    * `release_desc` - Configuration item release description.Note: This field may return null, which means no valid value was found.
    * `release_time` - Configuration item release time.Note: This field may return null, which means no valid value was found.
  * `file_config_release_list` - File configuration item release list.Note: This field may return null, which means no valid value was found.
    * `cluster_id` - Configuration item release cluster ID.Note: This field may return null, which means no valid value was found.
    * `cluster_name` - Configuration item release cluster name.Note: This field may return null, which means no valid value was found.
    * `config_id` - Configuration item release config ID.Note: This field may return null, which means no valid value was found.
    * `config_name` - Configuration item release config name.Note: This field may return null, which means no valid value was found.
    * `config_release_id` - Configuration item release ID.Note: This field may return null, which means no valid value was found.
    * `config_version` - Configuration item release config version.Note: This field may return null, which means no valid value was found.
    * `group_id` - Configuration item release config group ID.Note: This field may return null, which means no valid value was found.
    * `group_name` - Configuration item release config group name.Note: This field may return null, which means no valid value was found.
    * `namespace_id` - Configuration item release namespace ID.Note: This field may return null, which means no valid value was found.
    * `namespace_name` - Configuration item release namespace name.Note: This field may return null, which means no valid value was found.
    * `release_desc` - Configuration item release description.Note: This field may return null, which means no valid value was found.
    * `release_time` - Configuration item release time.Note: This field may return null, which means no valid value was found.
  * `package_id` - Package Id.Note: This field may return null, which means no valid value was found.
  * `package_name` - Package name.Note: This field may return null, which means no valid value was found.
  * `package_version` - Package version.Note: This field may return null, which means no valid value was found.
  * `public_config_release_list` - Release public config list.
    * `application_id` - Configuration item release application ID.Note: This field may return null, which means no valid value was found.
    * `cluster_id` - Configuration item release cluster ID.Note: This field may return null, which means no valid value was found.
    * `cluster_name` - Configuration item release cluster name.Note: This field may return null, which means no valid value was found.
    * `config_id` - Configuration item  ID.Note: This field may return null, which means no valid value was found.
    * `config_name` - Configuration item name.Note: This field may return null, which means no valid value was found.
    * `config_release_id` - Configuration item release ID.Note: This field may return null, which means no valid value was found.
    * `config_version` - Configuration version.Note: This field may return null, which means no valid value was found.
    * `group_id` - Configuration item release group ID.Note: This field may return null, which means no valid value was found.
    * `group_name` - Configuration item release group name.Note: This field may return null, which means no valid value was found.
    * `namespace_id` - Configuration item release namespace ID.Note: This field may return null, which means no valid value was found.
    * `namespace_name` - Configuration item release namespace name.Note: This field may return null, which means no valid value was found.
    * `release_desc` - Configuration item release description.Note: This field may return null, which means no valid value was found.
    * `release_time` - Configuration item release time.Note: This field may return null, which means no valid value was found.
  * `repo_name` - image name.Note: This field may return null, which means no valid value was found.
  * `tag_name` - image tag name.Note: This field may return null, which means no valid value was found.


