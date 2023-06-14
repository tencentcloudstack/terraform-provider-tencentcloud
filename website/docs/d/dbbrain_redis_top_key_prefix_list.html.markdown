---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_redis_top_key_prefix_list"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_redis_top_key_prefix_list"
description: |-
  Use this data source to query detailed information of dbbrain redis_top_key_prefix_list
---

# tencentcloud_dbbrain_redis_top_key_prefix_list

Use this data source to query detailed information of dbbrain redis_top_key_prefix_list

## Example Usage

```hcl
data "tencentcloud_dbbrain_redis_top_key_prefix_list" "redis_top_key_prefix_list" {
  instance_id = local.redis_id
  date        = "%s"
  product     = "redis"
}
```

## Argument Reference

The following arguments are supported:

* `date` - (Required, String) Query date, such as 2021-05-27, the earliest date can be the previous 30 days.
* `instance_id` - (Required, String) instance id.
* `product` - (Required, String) Service product type, supported values include `redis` - cloud database Redis.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - list of top key prefixes.
  * `ave_element_size` - Average element length.
  * `count` - The number of keys.
  * `item_count` - number of elements.
  * `key_pre_index` - key prefix.
  * `length` - Total occupied memory (Byte).
  * `max_element_size` - Maximum element length.


