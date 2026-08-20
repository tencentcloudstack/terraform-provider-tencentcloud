---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_edge_kv_list"
sidebar_current: "docs-tencentcloud-datasource-teo_edge_kv_list"
description: |-
  Use this data source to query the list of key names under a specified TEO EdgeKV namespace.
---

# tencentcloud_teo_edge_kv_list

Use this data source to query the list of key names under a specified TEO EdgeKV namespace.

## Example Usage

### Query all keys in a namespace

```hcl
data "tencentcloud_teo_edge_kv_list" "example" {
  zone_id   = "zone-2qtuhspy7cr6"
  namespace = "default"
}
```

### Query keys with a prefix filter

```hcl
data "tencentcloud_teo_edge_kv_list" "example" {
  zone_id   = "zone-2qtuhspy7cr6"
  namespace = "default"
  prefix    = "user_"
}
```

### Continue traversal with a cursor

```hcl
data "tencentcloud_teo_edge_kv_list" "example" {
  zone_id   = "zone-2qtuhspy7cr6"
  namespace = "default"
  cursor    = "next-page-cursor-from-previous-query"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String) EdgeKV namespace name.
* `zone_id` - (Required, String) Site ID.
* `cursor` - (Optional, String) Cursor position for traversal. Do not set this field for the first query; for subsequent queries, set it to the cursor returned by the previous query. After Read completes, this field is populated with the cursor from the last API response.
* `prefix` - (Optional, String) Key name prefix filter. Only keys starting with the specified prefix are returned. If not set, all keys are returned.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `keys` - List of key names in the specified namespace.


