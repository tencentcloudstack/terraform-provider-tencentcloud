---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_repositories"
sidebar_current: "docs-tencentcloud-datasource-tcr_repositories"
description: |-
  Use this data source to query detailed information of TCR repositories.
---

# tencentcloud_tcr_repositories

Use this data source to query detailed information of TCR repositories.

## Example Usage

```hcl
data "tencentcloud_tcr_repositories" "name" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) ID of the TCR instance that the repository belongs to.
* `namespace_name` - (Required, String) Name of the namespace that the repository belongs to.
* `repository_name` - (Optional, String) ID of the TCR repositories to query.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `repository_list` - Information list of the dedicated TCR repositories.
  * `brief_desc` - Brief description of the repository.
  * `create_time` - Create time.
  * `description` - Description of the repository.
  * `is_public` - Indicate that the repository is public or not.
  * `name` - Name of repository.
  * `namespace_name` - Name of the namespace that the repository belongs to.
  * `update_time` - Last update time.
  * `url` - URL of the repository.


