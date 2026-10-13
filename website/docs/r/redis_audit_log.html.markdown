---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_audit_log"
sidebar_current: "docs-tencentcloud-resource-redis_audit_log"
description: |-
  Provides a resource to manage Redis audit log configuration.
---

# tencentcloud_redis_audit_log

Provides a resource to manage Redis audit log configuration.

## Example Usage

```hcl
resource "tencentcloud_redis_audit_log" "example" {
  instance_id         = "crs-6eqwe3lt"
  log_sub_type        = "all"
  log_expire_day      = 7
  high_log_expire_day = 7
  degrade_strategy    = 500
}
```

## Argument Reference

The following arguments are supported:

* `high_log_expire_day` - (Required, Int) High-frequency log retention period in days. Valid values: `7` (7 days).
* `instance_id` - (Required, String, ForceNew) Instance ID, such as: crs-xjhsdj****, which can be copied from the instance list in the Redis console.
* `log_expire_day` - (Required, Int) Log retention period in days. Valid values: `7` (7 days), `30` (30 days).
* `log_sub_type` - (Required, String) Log sub-type. Valid values: `write` (write commands), `read` (read commands), `all` (read and write commands).
* `degrade_strategy` - (Optional, Int) Degradation strategy threshold in milliseconds. When the instance P99 latency reaches this threshold, audit logs will be automatically discarded to ensure business availability. Value range: [300, 1000]. Default value: `500`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeLogInstanceList](https://cloud.tencent.com/document/api/239/129828) - 查询日志实例列表
* Create: [OpenLog](https://cloud.tencent.com/document/api/239/129825) - 开启审计日志
* Read: [DescribeLogInstanceList](https://cloud.tencent.com/document/api/239/129828) - 查询日志实例列表
* Update: [DescribeLogInstanceList](https://cloud.tencent.com/document/api/239/129828) - 查询日志实例列表
* Update: [ModifyLog](https://cloud.tencent.com/document/api/239/129826) - 修改审计日志
* Delete: [CloseLog](https://cloud.tencent.com/document/api/239/129832) - 关闭审计日志
* Delete: [DescribeLogInstanceList](https://cloud.tencent.com/document/api/239/129828) - 查询日志实例列表

## Import

Redis audit log can be imported using the instance id, e.g.

```
terraform import tencentcloud_redis_audit_log.example crs-6eqwe3lt
```

