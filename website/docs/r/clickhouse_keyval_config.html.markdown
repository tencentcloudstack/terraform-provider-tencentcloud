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

* Create: [DescribeInstanceState](https://cloud.tencent.com/document/api/1299/89676) - 获取集群实例状态
* Create: [ModifyInstanceKeyValConfigs](https://cloud.tencent.com/document/api/1299/86159) - KV模式修改配置接口
* Read: [DescribeInstanceKeyValConfigs](https://cloud.tencent.com/document/api/1299/86160) - 在集群详情页面获取所有参数列表
* Update: [DescribeInstanceState](https://cloud.tencent.com/document/api/1299/89676) - 获取集群实例状态
* Update: [ModifyInstanceKeyValConfigs](https://cloud.tencent.com/document/api/1299/86159) - KV模式修改配置接口
* Delete: [DescribeInstanceState](https://cloud.tencent.com/document/api/1299/89676) - 获取集群实例状态
* Delete: [ModifyInstanceKeyValConfigs](https://cloud.tencent.com/document/api/1299/86159) - KV模式修改配置接口

## Import

clickhouse config can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_keyval_config.config cdwch-datuhk3z#max_open_files#50000
```

