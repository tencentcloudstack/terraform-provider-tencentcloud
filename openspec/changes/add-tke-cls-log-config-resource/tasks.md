# Tasks: Add Update Support for TKE CLS Log Config Resource

## Implementation

- [ ] 修改 `tencentcloud/services/tke/resource_tc_kubernetes_log_config.go`
  - [ ] 在 Schema 中移除 `log_config_name` 的 `ForceNew: true`
  - [ ] 移除 `log_config` 的 `ForceNew: true`（支持更新配置内容）
  - [ ] 添加 `Update` 函数：`resourceTencentCloudKubernetesLogConfigUpdate`
  - [ ] 注册 Update 函数到 Resource 结构

- [ ] 修改 `tencentcloud/services/tke/service_tencentcloud_tke.go`
  - [ ] 添加 `ModifyKubernetesLogConfig` 方法
  - [ ] 处理 ModifyLogConfig API 调用和错误

## Testing

- [ ] 验证编译: `go build ./...`
- [ ] 运行 linter: `make lint`
- [ ] 格式化代码: `make fmt`

## Documentation

- [ ] 更新 `tencentcloud/services/tke/resource_tc_kubernetes_log_config.md`
  - [ ] 添加 Update 用法说明
  - [ ] 添加更新示例