---
subcategory: "TencentCloud Automation Tools(TAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tat_invoker_config"
sidebar_current: "docs-tencentcloud-resource-tat_invoker_config"
description: |-
  Provides a resource to create a tat invoker_config
---

# tencentcloud_tat_invoker_config

Provides a resource to create a tat invoker_config

## Example Usage

```hcl
resource "tencentcloud_tat_invoker_config" "invoker_config" {
  invoker_id     = "ivk-cas4upyf"
  invoker_status = "on"
}
```

## Argument Reference

The following arguments are supported:

* `invoker_id` - (Required, String) ID of the invoker to be enabled.
* `invoker_status` - (Required, String) Invoker on and off state, Values: `on`, `off`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DisableInvoker](https://cloud.tencent.com/document/api/1340/61758) - 停用执行器
* Create: [EnableInvoker](https://cloud.tencent.com/document/api/1340/61757) - 启用执行器
* Read: [DescribeInvokers](https://cloud.tencent.com/document/api/1340/61759) - 查询执行器
* Update: [DisableInvoker](https://cloud.tencent.com/document/api/1340/61758) - 停用执行器
* Update: [EnableInvoker](https://cloud.tencent.com/document/api/1340/61757) - 启用执行器

## Import

tat invoker_config can be imported using the id, e.g.

```
terraform import tencentcloud_tat_invoker_config.invoker_config invoker_config_id
```

