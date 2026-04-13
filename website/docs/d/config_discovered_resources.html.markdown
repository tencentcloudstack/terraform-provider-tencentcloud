---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_discovered_resources"
sidebar_current: "docs-tencentcloud-datasource-config_discovered_resources"
description: |-
  Use this data source to query detailed information of Config discovered resources.
---

# tencentcloud_config_discovered_resources

Use this data source to query detailed information of Config discovered resources.

## Example Usage

### Query all discovered resources

```hcl
data "tencentcloud_config_discovered_resources" "example" {}
```

### Query by resource ID filter

```hcl
data "tencentcloud_config_discovered_resources" "example" {
  filters {
    name   = "resourceId"
    values = ["ins-pbu2hghz"]
  }
  order_type = "desc"
}
```

### Query by resource name and tags

```hcl
data "tencentcloud_config_discovered_resources" "example" {
  filters {
    name   = "resourceName"
    values = ["my-cvm"]
  }

  tags {
    tag_key   = "env"
    tag_value = "prod"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. Supported filter names: resourceName (resource name), resourceId (resource ID).
* `order_type` - (Optional, String) Sort type. Valid values: asc, desc.
* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, List) Tag filter conditions.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name. Valid values: resourceName, resourceId.
* `values` - (Required, List) Filter field values.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `resource_list` - Discovered resource list.
  * `compliance_result` - Compliance result. Valid values: COMPLIANT, NON_COMPLIANT.
  * `resource_create_time` - Resource creation time.
  * `resource_delete` - Resource deletion mark. Valid values: 1 (deleted), 2 (not deleted).
  * `resource_id` - Resource ID.
  * `resource_name` - Resource name.
  * `resource_region` - Resource region.
  * `resource_status` - Resource status.
  * `resource_type` - Resource type.
  * `resource_zone` - Resource availability zone.
  * `tags` - Resource tag list.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.


