# Design: Add Update Support for TKE CLS Log Config Resource

## Overview

为 `tencentcloud_kubernetes_log_config` 资源添加 Update 功能，使其支持通过 `ModifyLogConfig` API 更新日志采集配置。

## API Details

### ModifyLogConfig 接口

| 参数名 | 必选 | 类型 | 描述 |
|--------|------|------|------|
| ClusterId | 是 | String | 集群 ID |
| LogConfig | 是 | String | 日志采集配置的 JSON 表达 |
| ClusterType | 否 | String | 集群类型，支持 `tke`（标准集群）、`eks`（serverless 集群）|

### LogConfig JSON 结构

参考: https://cloud.tencent.com/document/product/457/48425

```json
{
  "apiVersion": "cls.cloud.tencent.com/v1",
  "kind": "LogConfig",
  "metadata": {
    "name": "config-name"
  },
  "spec": {
    "clsDetail": {},
    "inputDetail": {},
    "kafkaDetail": {}
  }
}
```

## Implementation Design

### 1. Resource Schema 调整

现有 schema 字段：

| 字段 | 类型 | Required | ForceNew | 说明 |
|------|------|----------|----------|------|
| log_config | String | Yes | Yes | 日志采集配置的 JSON 表达 |
| log_config_name | String | Yes | Yes | 日志配置名称 |
| cluster_id | String | Yes | Yes | 集群 ID |
| logset_id | String | No | Yes | CLS 日志集 ID |
| cluster_type | String | No | Yes | 集群类型 |

**调整**：移除 `log_config_name` 的 `ForceNew`，使其支持更新时修改名称。

### 2. Update 函数实现

```go
func resourceTencentCloudKubernetesLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
    // 1. 解析 ID 获取 clusterId 和 logConfigName
    // 2. 检查变更字段
    // 3. 调用 ModifyLogConfig API
    // 4. 等待配置同步
    // 5. 调用 Read 函数刷新状态
}
```

### 3. Service 层更新

在 `TkeService` 中添加 `ModifyKubernetesLogConfig` 方法：

```go
func (me *TkeService) ModifyKubernetesLogConfig(ctx context.Context, clusterId, logConfig, clusterType string) error {
    request := tkev20180525.NewModifyLogConfigRequest()
    request.ClusterId = &clusterId
    request.LogConfig = &logConfig
    request.ClusterType = &clusterType

    _, err := me.client.UseTkeV20180525Client().ModifyLogConfigWithContext(ctx, request)
    return err
}
```

### 4. 错误处理

- 重试机制：使用 `WriteRetryTimeout` 进行最终一致性重试
- 错误码处理：处理 `FailedOperation` 相关错误

## Dependencies

- TKE API v20180525 SDK
- 现有的 Read/Delete 函数（可复用）

## Testing Strategy

1. 单元测试：验证 Update 函数逻辑
2. 集成测试：创建 -> 更新 -> 验证状态