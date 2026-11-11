## Why

CLS (日志服务) 提供了指标订阅 (Metric Subscribe) 功能，允许用户将日志主题中的指标数据订阅并采集到云监控。当前 Terraform Provider 中没有 `tencentcloud_cls_metric_subscribe` 资源，用户无法通过 Terraform 管理指标订阅配置的完整生命周期 (创建、读取、更新、删除)，只能通过控制台或 API 手动操作，导致基础设施漂移且无法统一管理。

## What Changes

- 新增 RESOURCE_KIND_GENERAL 资源 `tencentcloud_cls_metric_subscribe`，支持完整的 CRUD 操作
- 资源使用 CLS API v20201016 的以下接口:
  - `CreateMetricSubscribe` — 创建指标订阅配置
  - `DescribeMetricSubscribes` — 查询指标订阅配置 (用于 Read)
  - `ModifyMetricSubscribe` — 修改指标订阅配置 (用于 Update)
  - `DeleteMetricSubscribe` — 删除指标订阅配置
- 资源使用复合 ID 格式: `topicId#taskId` (由 `tccommon.FILED_SP` 分隔)
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源
- 生成对应的单元测试文件 (使用 gomonkey mock 云 API)
- 生成对应的 `.md` 文档

## Capabilities

### New Capabilities
- `cls-metric-subscribe`: CLS 指标订阅资源配置管理，包括名称、日志主题 ID、云产品命名空间、指标配置信息、实例配置信息、任务开关等参数的增删改查

### Modified Capabilities
<!-- 无修改的已有能力 -->

## Impact

- **新增文件**:
  - `tencentcloud/services/cls/resource_tc_cls_metric_subscribe.go` — 资源 CRUD 实现
  - `tencentcloud/services/cls/resource_tc_cls_metric_subscribe.md` — 资源文档
  - `tencentcloud/services/cls/resource_tc_cls_metric_subscribe_test.go` — 单元测试
- **修改文件**:
  - `tencentcloud/services/cls/service_tencentcloud_cls.go` — 新增 `DescribeClsMetricSubscribeById` 方法
  - `tencentcloud/provider.go` — 注册 `tencentcloud_cls_metric_subscribe` 资源
  - `tencentcloud/provider.md` — 添加资源文档条目
- **API 依赖**: CLS API v20201016 (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`)
- **兼容性**: 向后兼容，纯新增资源，不影响现有资源
