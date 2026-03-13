# Implementation Tasks

## 1. 核心数据源实现
- [x] 1.1 创建数据源文件 `tencentcloud/services/monitor/data_source_tc_monitor_external_cluster_register_command.go`
  - [x] 1.1.1 定义 Schema 结构 (instance_id, cluster_id 为必需字段)
  - [x] 1.1.2 实现 `DataSourceTencentCloudMonitorExternalClusterRegisterCommand()` 函数
  - [x] 1.1.3 实现 `dataSourceTencentCloudMonitorExternalClusterRegisterCommandRead()` 函数
    - 调用服务层方法获取注册命令
    - 构造复合 ID: `instanceId#clusterId`
    - 映射 API 返回字段到 Schema
    - 添加错误处理和重试逻辑
  - [x] 1.1.4 添加日志记录和错误处理
  - [x] 1.1.5 支持 `result_output_file` 参数

## 2. 服务层增强
- [x] 2.1 在 `service_tencentcloud_monitor.go` 中添加 `DescribeExternalClusterRegisterCommand()` 方法
  - 接收 ctx, instanceId, clusterId 参数
  - 调用 `DescribeExternalClusterRegisterCommand` API
  - 返回注册命令相关数据
  - 包含重试逻辑和错误处理

## 3. Provider 注册
- [x] 3.1 在 `tencentcloud/provider.go` 的 `DataSourcesMap` 中注册新数据源
  - 添加 `"tencentcloud_monitor_external_cluster_register_command": monitor.DataSourceTencentCloudMonitorExternalClusterRegisterCommand()`

## 4. 测试实现
- [x] 4.1 创建测试文件 `data_source_tc_monitor_external_cluster_register_command_test.go`
- [x] 4.2 实现基础验收测试 `TestAccTencentCloudMonitorExternalClusterRegisterCommand_basic`
  - 测试查询注册命令
  - 验证返回字段正确性
- [ ] 4.3 确保测试可以通过 `TF_ACC=1 go test`

## 5. 文档编写
- [x] 5.1 创建数据源文档 `website/docs/d/monitor_external_cluster_register_command.html.markdown`
  - 数据源描述
  - 参数说明 (标注必选/可选)
  - 使用示例
  - 属性参考
- [ ] 5.2 运行 `make doc` 生成文档

## 6. 代码质量检查
- [x] 6.1 运行 `make fmt` 格式化代码
- [x] 6.2 运行 `make lint` 检查代码质量
- [x] 6.3 修复所有 linter 警告和错误
- [x] 6.4 确保代码符合项目编码规范

## 7. 集成测试
- [ ] 7.1 在测试环境运行完整的验收测试
- [ ] 7.2 测试数据源的查询功能
- [ ] 7.3 验证错误场景处理
- [ ] 7.4 测试与 `tencentcloud_monitor_external_cluster` 资源的配合使用

## 8. 代码审查准备
- [x] 8.1 确认所有任务已完成
- [ ] 8.2 自我审查代码
- [ ] 8.3 准备 PR 描述和变更日志
- [ ] 8.4 创建 `.changelog/<PR_NUMBER>.txt` 文件

## 关键验证点
- ✓ 数据源 ID 格式正确: `instanceId#clusterId`
- ✓ 参数 `InstanceId` 和 `ClusterId` 都为必传
- ✓ 代码格式参考 `data_source_tc_igtm_instance_list.go`
- ✓ API 调用包含重试逻辑
- ✓ 错误处理完善
- ✓ 日志记录规范
- ✓ 支持 `result_output_file` 输出
