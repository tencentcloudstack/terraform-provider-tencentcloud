---
subcategory: "CdwDoris"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwdoris_instances"
sidebar_current: "docs-tencentcloud-datasource-cdwdoris_instances"
description: |-
  Use this data source to query detailed information of cdwdoris instances
---

# tencentcloud_cdwdoris_instances

Use this data source to query detailed information of cdwdoris instances

## Example Usage

### Query all cdwdoris instances

```hcl
data "tencentcloud_cdwdoris_instances" "example" {}
```

### Query cdwdoris instances by filter

```hcl
# by instance Id
data "tencentcloud_cdwdoris_instances" "example" {
  search_instance_id = "cdwdoris-rhbflamd"
}

# by instance name
data "tencentcloud_cdwdoris_instances" "example" {
  search_instance_name = "tf-example"
}

# by instance tags
data "tencentcloud_cdwdoris_instances" "example" {
  search_tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
    all_value = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `search_instance_id` - (Optional, String) The name of the cluster ID for the search.
* `search_instance_name` - (Optional, String) The cluster name for the search.
* `search_tags` - (Optional, List) Search tag list.

The `search_tags` object supports the following:

* `all_value` - (Optional, Int) 1 means only the tag key is entered without a value, and 0 means both the key and the value are entered.
* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances_list` - Quantities of instances array.


