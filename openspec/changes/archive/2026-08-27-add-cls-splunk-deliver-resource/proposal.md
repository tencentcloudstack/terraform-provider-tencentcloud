## Why

腾讯云 CLS（日志服务）支持将日志数据投递到 Splunk，但当前 Terraform Provider 缺少对 Splunk 投递任务的管理能力。用户无法通过 IaC 方式自动化创建、配置和管理 Splunk 投递任务，需要手动在控制台操作，增加了运维复杂度和出错风险。新增此资源可以让用户通过 Terraform 统一管理 CLS 日志到 Splunk 的投递链路。

## What Changes

- 新增 `tencentcloud_cls_splunk_deliver` Terraform 资源，支持 Splunk 投递任务的完整 CRUD 生命周期管理
- 通过 `CreateSplunkDeliver` API 创建 Splunk 投递任务
- 通过 `DescribeSplunkDelivers` API 查询 Splunk 投递任务详情
- 通过 `ModifySplunkDeliver` API 修改 Splunk 投递任务配置
- 通过 `DeleteSplunkDeliver` API 删除 Splunk 投递任务
- 支持资源导入（import），通过 `task_id` 和 `topic_id` 联合 ID 导入已有投递任务

## Capabilities

### New Capabilities
- `cls-splunk-deliver-resource`: 提供 CLS Splunk 投递任务的 Terraform 资源管理能力，包括创建、读取、更新、删除和导入操作

### Modified Capabilities
<!-- No existing capabilities are modified -->

## Impact

- 新增文件: `tencentcloud/services/cls/resource_tc_cls_splunk_deliver.go`
- 新增文件: `tencentcloud/services/cls/resource_tc_cls_splunk_deliver_test.go`
- 新增文件: `tencentcloud/services/cls/resource_tc_cls_splunk_deliver.md`
- 修改文件: `tencentcloud/provider.go` - 注册新资源
- 修改文件: `tencentcloud/provider.md` - 添加资源文档入口
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016` (已存在于 vendor 中)