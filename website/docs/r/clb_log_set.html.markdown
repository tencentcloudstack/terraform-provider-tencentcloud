---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_log_set"
sidebar_current: "docs-tencentcloud-resource-clb_log_set"
description: |-
  Provides a resource to create an exclusive CLB Logset.
---

# tencentcloud_clb_log_set

Provides a resource to create an exclusive CLB Logset.

## Example Usage

```hcl
resource "tencentcloud_clb_log_set" "foo" {
  period = 7
}
```

## Argument Reference

The following arguments are supported:

* `period` - (Optional, Int, ForceNew) Logset retention period in days. Maximun value is `90`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Logset creation time.
* `name` - Logset name, which unique and fixed `clb_logset` among all CLS logsets.
* `topic_count` - Number of log topics in logset.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateClsLogSet](https://cloud.tencent.com/document/api/214/46803) - This API is used to create a CLB exclusive logset for storing CLB logs.
* Create: [DescribeClsLogSet](https://cloud.tencent.com/document/api/214/46801) - This API is used to get the CLB exclusive logset.
* Read: [DescribeLogsets](https://cloud.tencent.com/document/api/614/58624) - This API is used to get the list of logsets.
* Delete: [DeleteLogset](https://cloud.tencent.com/document/api/614/58625) - This API is used to delete a logset.

## Import

CLB log set can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_logset.foo 4eb9e3a8-9c42-4b32-9ddf-e215e9c92764
```

