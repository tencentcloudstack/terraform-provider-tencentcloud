---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_db_parameters"
sidebar_current: "docs-tencentcloud-resource-dcdb_db_parameters"
description: |-
  Provides a resource to create a dcdb db_parameters
---

# tencentcloud_dcdb_db_parameters

Provides a resource to create a dcdb db_parameters

## Example Usage

```hcl
resource "tencentcloud_dcdb_db_parameters" "db_parameters" {
  instance_id = "%s"
  params {
    param = "max_connections"
    value = "9999"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of instance.
* `params` - (Required, List) Parameter list, each element is a combination of Param and Value.

The `params` object supports the following:

* `param` - (Required, String) The name of parameter.
* `value` - (Required, String) The value of parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeDBParameters](https://cloud.tencent.com/document/api/557/19992) - 查看数据库参数
* Create: [ModifyDBParameters](https://cloud.tencent.com/document/api/557/19982) - 修改数据库参数
* Read: [DescribeDBParameters](https://cloud.tencent.com/document/api/557/19992) - 查看数据库参数
* Update: [DescribeDBParameters](https://cloud.tencent.com/document/api/557/19992) - 查看数据库参数
* Update: [ModifyDBParameters](https://cloud.tencent.com/document/api/557/19982) - 修改数据库参数

## Import

dcdb db_parameters can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_db_parameters.db_parameters instanceId#paramName
```

