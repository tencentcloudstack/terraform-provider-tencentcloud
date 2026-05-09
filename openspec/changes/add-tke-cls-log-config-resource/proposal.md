# Change: Add TKE CLS Log Config Resource

## Why

用户需要在 Terraform 中管理 TKE 集群的 CLS 日志采集配置。当前 TKE 服务的 `tencentcloud_kubernetes_log_config` 资源缺少 Update 功能，只能创建和删除，无法修改已有的日志采集配置。这限制了用户在 Terraform 中完整管理日志配置的生命周期。

具体问题：
1. 无法修改现有的日志采集配置
2. 日志配置变更需要通过其他方式（如 kubectl）手动操作
3. 无法通过 Terraform 统一管理日志配置

## What Changes

- 为 `tencentcloud_kubernetes_log_config` 资源添加 Update 功能
- 支持通过 `ModifyLogConfig` 接口更新日志采集配置
- 保持现有 Create/Delete/Read 功能不变
- 资源 ID 格式保持为 `ClusterId#name`

## Impact

- **修改范围**: TKE 服务 `tencentcloud_kubernetes_log_config` 资源
- **受影响的服务**: TKE (tencentcloud/services/tke)
- **修改文件**:
  - `tencentcloud/services/tke/resource_tc_kubernetes_log_config.go`
  - `tencentcloud/services/tke/service_tencentcloud_tke.go`
- **API 依赖**:
  - TKE API v20180525: `ModifyLogConfig`
  - 文档: https://cloud.tencent.com/document/api/457/130677
- **兼容性**: 向后兼容，无破坏性变更