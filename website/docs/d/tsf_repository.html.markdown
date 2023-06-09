---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_repository"
sidebar_current: "docs-tencentcloud-datasource-tsf_repository"
description: |-
  Use this data source to query detailed information of tsf repository
---

# tencentcloud_tsf_repository

Use this data source to query detailed information of tsf repository

## Example Usage

```hcl
data "tencentcloud_tsf_repository" "repository" {
  search_word     = "test"
  repository_type = "default"
}
```

## Argument Reference

The following arguments are supported:

* `repository_type` - (Optional, String) Repository type (default Repository: default, private Repository: private).
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) Query keywords (search by Repository name).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - A list of Repository information that meets the query criteria.
  * `content` - Repository information list. Note: This field may return null, indicating that no valid value can be obtained.
    * `bucket_name` - Repository bucket name. Note: This field may return null, indicating that no valid value can be obtained.
    * `bucket_region` - Repository region. Note: This field may return null, indicating that no valid value can be obtained.
    * `create_time` - CreationTime. Note: This field may return null, indicating that no valid values can be obtained.
    * `directory` - Repository Directory. Note: This field may return null, indicating that no valid value can be obtained.
    * `is_used` - Whether the repository is being used. Note: This field may return null, indicating that no valid value can be obtained.
    * `repository_desc` - Repository description (default warehouse: default, private warehouse: private).
    * `repository_id` - repository Id.
    * `repository_name` - Repository Name.
    * `repository_type` - Repository type (default Repository: default, private Repository: private).
  * `total_count` - Total Repository.


