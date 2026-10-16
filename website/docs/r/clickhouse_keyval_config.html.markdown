---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_keyval_config"
sidebar_current: "docs-tencentcloud-resource-clickhouse_keyval_config"
description: |-
  Provides a resource to create a clickhouse keyval_config
---

# tencentcloud_clickhouse_keyval_config

Provides a resource to create a clickhouse keyval_config

## Example Usage

```hcl
resource "tencentcloud_clickhouse_keyval_config" "keyval_config" {
  instance_id = "cdwch-datuhk3z"
  items {
    conf_key   = "max_open_files"
    conf_value = "50000"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `items` - (Required, List) configuration list.

The `items` object supports the following:

* `conf_key` - (Required, String) Instance config key.
* `conf_value` - (Required, String) Instance config value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeInstanceState](https://www.tencentcloud.com/search/?wd=DescribeInstanceState&type=document&lang=en)
* Create: [ModifyInstanceKeyValConfigs](https://www.tencentcloud.com/search/?wd=ModifyInstanceKeyValConfigs&type=document&lang=en)
* Read: [DescribeInstanceKeyValConfigs](https://www.tencentcloud.com/search/?wd=DescribeInstanceKeyValConfigs&type=document&lang=en)
* Update: [DescribeInstanceState](https://www.tencentcloud.com/search/?wd=DescribeInstanceState&type=document&lang=en)
* Update: [ModifyInstanceKeyValConfigs](https://www.tencentcloud.com/search/?wd=ModifyInstanceKeyValConfigs&type=document&lang=en)
* Delete: [DescribeInstanceState](https://www.tencentcloud.com/search/?wd=DescribeInstanceState&type=document&lang=en)
* Delete: [ModifyInstanceKeyValConfigs](https://www.tencentcloud.com/search/?wd=ModifyInstanceKeyValConfigs&type=document&lang=en)

## Import

clickhouse config can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_keyval_config.config cdwch-datuhk3z#max_open_files#50000
```

