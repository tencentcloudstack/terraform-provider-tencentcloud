---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_redis_top_big_keys"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_redis_top_big_keys"
description: |-
  Use this data source to query detailed information of dbbrain redis_top_big_keys
---

# tencentcloud_dbbrain_redis_top_big_keys

Use this data source to query detailed information of dbbrain redis_top_big_keys

## Example Usage

```hcl
data "tencentcloud_dbbrain_redis_top_big_keys" "redis_top_big_keys" {
  instance_id = local.redis_id
  date        = "%s"
  product     = "redis"
  sort_by     = "Capacity"
  key_type    = "string"
}
```

## Argument Reference

The following arguments are supported:

* `date` - (Required, String) Query date, such as 2021-05-27, the earliest date can be the previous 30 days.
* `instance_id` - (Required, String) instance id.
* `product` - (Required, String) Service product type, supported values include `redis` - cloud database Redis.
* `key_type` - (Optional, String) Key type filter condition, the default is no filter, the value includes `string`, `list`, `set`, `hash`, `sortedset`, `stream`.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) Sorting field, the value includes `Capacity` - memory, `ItemCount` - number of elements, the default is `Capacity`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `top_keys` - list of top keys.
  * `encoding` - key encoding method.
  * `expire_time` - Key expiration timestamp (in milliseconds), 0 means no expiration time is set.
  * `item_count` - number of elements.
  * `key` - key name.
  * `length` - Key memory size, unit Byte.
  * `max_element_size` - Maximum element length.
  * `type` - key type.


