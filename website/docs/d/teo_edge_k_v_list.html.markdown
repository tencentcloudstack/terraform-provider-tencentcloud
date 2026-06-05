---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_edge_k_v_list"
sidebar_current: "docs-tencentcloud-datasource-teo_edge_k_v_list"
description: |-
  Use this data source to query detailed information of TencentCloud EdgeOne (TEO) Edge KV key names list
---

# tencentcloud_teo_edge_k_v_list

Use this data source to query detailed information of TencentCloud EdgeOne (TEO) Edge KV key names list

## Example Usage

### Query all keys in a namespace

```hcl
data "tencentcloud_teo_edge_k_v_list" "example" {
  zone_id   = "zone-2q1ysez95gao"
  namespace = "my_namespace"
}
```

### Query keys with prefix filter

```hcl
data "tencentcloud_teo_edge_k_v_list" "example" {
  zone_id   = "zone-2q1ysez95gao"
  namespace = "my_namespace"
  prefix    = "config_"
}
```

### Query keys with cursor for pagination

```hcl
data "tencentcloud_teo_edge_k_v_list" "example" {
  zone_id   = "zone-2q1ysez95gao"
  namespace = "my_namespace"
  cursor    = "eyJjIjoiTVRBd01EQXdNREF3TURBd01EQXdNREE9IiwidCI6ImFhYSJ9"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String) Namespace name.
* `zone_id` - (Required, String) Zone ID.
* `cursor` - (Optional, String) Cursor position. Identifies the starting position of the current query for traversing large amounts of data. Leave empty for the first query to start from the beginning; for subsequent queries, fill in the Cursor value returned from the previous response to continue traversal.
* `prefix` - (Optional, String) Key name prefix filter. Only returns keys that start with the specified prefix, length 1-512 characters.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `keys` - List of key names.


