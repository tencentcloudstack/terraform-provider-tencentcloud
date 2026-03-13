# Implementation Tasks

## 1. 核心资源实现
- [x] 1.1 创建资源文件 `tencentcloud/services/monitor/resource_tc_monitor_external_cluster.go`
  - [x] 1.1.1 定义 Schema 结构 (包含所有必需和可选字段)
  - [x] 1.1.2 实现 `ResourceTencentCloudMonitorExternalCluster()` 函数
  - [x] 1.1.3 实现 `resourceTencentCloudMonitorExternalClusterCreate()` 函数
    - 调用 `CreateExternalCluster` API
    - 构造复合 ID: `instanceId#clusterId`
    - 添加错误处理和重试逻辑
  - [x] 1.1.4 实现 `resourceTencentCloudMonitorExternalClusterRead()` 函数
    - 解析复合 ID
    - 调用 `DescribePrometheusClusterAgents` API
    - 注意 `ClusterIds` 参数格式为字符串数组
    - 正确映射 `ClusterType` 字段到状态
  - [x] 1.1.5 实现 `resourceTencentCloudMonitorExternalClusterUpdate()` 函数
    - 如果需要更新则调用相应 API (检查是否有 Modify 接口)
    - 如果不支持更新，返回合适的错误信息
  - [x] 1.1.6 实现 `resourceTencentCloudMonitorExternalClusterDelete()` 函数
    - 解析复合 ID
    - 从状态中获取 `cluster_type` 字段
    - 构造 `Agents` 参数对象
    - 调用 `DeletePrometheusClusterAgent` API
  - [x] 1.1.7 添加 Import 支持
  - [x] 1.1.8 添加日志记录和错误处理

## 2. 服务层增强 (如果需要)
- [x] 2.1 检查 `service_tencentcloud_monitor.go` 是否需要新增辅助方法
- [x] 2.2 如需要，添加 `DescribeMonitorExternalClusterById()` 方法
- [x] 2.3 如需要,添加 `DeleteMonitorExternalCluster()` 方法

## 3. Provider 注册
- [x] 3.1 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册新资源
  - 添加 `"tencentcloud_monitor_external_cluster": monitor.ResourceTencentCloudMonitorExternalCluster()`

## 4. 测试实现
- [x] 4.1 创建测试文件 `resource_tc_monitor_external_cluster_test.go`
- [x] 4.2 实现基础验收测试 `TestAccTencentCloudMonitorExternalCluster_basic`
  - 测试创建外部集群
  - 测试读取集群信息
  - 测试删除集群
- [ ] 4.3 实现完整参数测试 (包含所有可选字段)
- [ ] 4.4 确保测试可以通过 `TF_ACC=1 go test`

## 5. 文档编写
- [x] 5.1 创建资源文档 `website/docs/r/monitor_external_cluster.html.markdown`
  - 资源描述
  - 参数说明 (标注必选/可选)
  - 使用示例 (基础示例 + 完整示例)
  - 属性参考
  - Import 说明
- [ ] 5.2 运行 `make doc` 生成文档

## 6. 代码质量检查
- [x] 6.1 运行 `make fmt` 格式化代码
- [x] 6.2 运行 `make lint` 检查代码质量
- [x] 6.3 修复所有 linter 警告和错误
- [x] 6.4 确保代码符合项目编码规范

## 7. 集成测试
- [ ] 7.1 在测试环境运行完整的验收测试
- [ ] 7.2 测试资源的完整生命周期 (创建 → 读取 → 更新 → 删除)
- [ ] 7.3 测试 Import 功能
- [ ] 7.4 验证错误场景处理

## 8. 代码审查准备
- [x] 8.1 确认所有任务已完成
- [ ] 8.2 自我审查代码
- [ ] 8.3 准备 PR 描述和变更日志
- [ ] 8.4 创建 `.changelog/<PR_NUMBER>.txt` 文件

## 关键验证点
- ✓ 资源 ID 格式正确: `instanceId#clusterId`
- ✓ `ClusterIds` 参数传递为字符串数组格式
- ✓ `ClusterType` 字段正确保存为 Computed 字段
- ✓ Delete 操作能正确从状态中获取 `cluster_type`
- ✓ 代码格式参考 `resource_tc_igtm_strategy.go`
- ✓ 所有 API 调用包含重试逻辑
- ✓ 错误处理完善
- ✓ 日志记录规范
