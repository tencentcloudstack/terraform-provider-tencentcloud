---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_enable_optimal_switching"
sidebar_current: "docs-tencentcloud-resource-css_enable_optimal_switching"
description: |-
  Provides a resource to create a css enable_optimal_switching
---

# tencentcloud_css_enable_optimal_switching

Provides a resource to create a css enable_optimal_switching

~> **NOTE:** This resource is only valid when the push stream. When the push stream ends, it will be deleted.

## Example Usage

```hcl
resource "tencentcloud_css_enable_optimal_switching" "enable_optimal_switching" {
  stream_name     = "1308919341_test"
  enable_switch   = 1
  host_group_name = "test-group"
}
```

## Argument Reference

The following arguments are supported:

* `stream_name` - (Required, String, ForceNew) Stream id.
* `enable_switch` - (Optional, Int) `0`:disabled, `1`:enable.
* `host_group_name` - (Optional, String, ForceNew) Group name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [EnableOptimalSwitching](https://cloud.tencent.com/document/api/267/100421) - 启用择优调度
* Read: [DescribeBackupStreamList](https://cloud.tencent.com/document/api/267/100422) - 查询直播中的主备流
* Update: [EnableOptimalSwitching](https://cloud.tencent.com/document/api/267/100421) - 启用择优调度

## Import

css domain can be imported using the id, e.g.

```
terraform import tencentcloud_css_enable_optimal_switching.enable_optimal_switching streamName
```

