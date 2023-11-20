---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_data_source_without_info"
sidebar_current: "docs-tencentcloud-datasource-wedata_data_source_without_info"
description: |-
  Use this data source to query detailed information of wedata data_source_without_info
---

# tencentcloud_wedata_data_source_without_info

Use this data source to query detailed information of wedata data_source_without_info

## Example Usage

```hcl
data "tencentcloud_wedata_data_source_without_info" "example" {
  filters {
    name   = "ownerProjectId"
    values = ["1612982498218618880"]
  }

  order_fields {
    name      = "create_time"
    direction = "DESC"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filters.
* `order_fields` - (Optional, List) OrderFields.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) Filter name.
* `values` - (Optional, Set) Filter value.

The `order_fields` object supports the following:

* `direction` - (Required, String) OrderFields rule.
* `name` - (Required, String) OrderFields name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Data.
  * `app_id` - Appid.
  * `author` - Has Author.
  * `authority_project_name` - Datasource AuthorityProjectName.
  * `authority_user_name` - Datasource AuthorityUserName.
  * `biz_params_string` - Biz params json string.
  * `biz_params` - Biz params.
  * `category` - Datasource category.
  * `cluster_id` - Datasource cluster id.
  * `cluster_name` - Datasource cluster name.
  * `create_time` - CreateTime.
  * `data_source_status` - DatasourceDataSourceStatus.
  * `database_name` - DatabaseName.
  * `deliver` - Can Deliver.
  * `description` - Description.
  * `display` - Datasource display name.
  * `edit` - Datasource can Edit.
  * `id` - ID.
  * `instance` - Instance.
  * `modified_time` - Datasource ModifiedTime.
  * `name` - Datasource name.
  * `owner_account_name` - Datasource owner account name.
  * `owner_account` - Datasource owner account.
  * `owner_project_id` - Datasource owner project id.
  * `owner_project_ident` - Datasource OwnerProjectIdent.
  * `owner_project_name` - Datasource OwnerProjectName.
  * `params_string` - Params json string.
  * `params` - Datasource params.
  * `region` - Datasource engin cluster region.
  * `show_type` - Datasource show type.
  * `status` - Datasource status.
  * `type` - Datasource type.


