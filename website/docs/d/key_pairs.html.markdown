---
subcategory: "CVM"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_key_pairs"
sidebar_current: "docs-tencentcloud-datasource-key_pairs"
description: |-
  Use this data source to query key pairs.
---

# tencentcloud_key_pairs

Use this data source to query key pairs.

## Example Usage

```hcl
data "tencentcloud_key_pairs" "foo" {
  key_id = "skey-ie97i3ml"
}

data "tencentcloud_key_pairs" "name" {
  key_name = "^test$"
}
```

## Argument Reference

The following arguments are supported:

* `key_id` - (Optional) ID of the key pair to be queried.
* `key_name` - (Optional) Name of the key pair to be queried. Support regular expression search, only `^` and `$` are supported.
* `project_id` - (Optional) Project id of the key pair to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `key_pair_list` - An information list of key pair. Each element contains the following attributes:
  * `create_time` - Creation time of the key pair.
  * `key_id` - ID of the key pair.
  * `key_name` - Name of the key pair.
  * `project_id` - Project id of the key pair.
  * `public_key` - public key of the key pair.


