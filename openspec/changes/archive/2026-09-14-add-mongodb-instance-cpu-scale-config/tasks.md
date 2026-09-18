# 实现任务清单

## 1. 服务层实现

- [x] 1.1 在 `service_tencentcloud_mongodb.go` 中添加 `DescribeMongodbInstanceById` 方法（调用 DescribeDBInstances API 查询单个实例）
- [x] 1.2 在 `service_tencentcloud_mongodb.go` 中添加 `ScaleUpMongodbDBInstanceCpu` 方法（调用 ScaleUpDBInstanceCpu API 开启 CPU 弹性扩容）
- [x] 1.3 在 `service_tencentcloud_mongodb.go` 中添加 `ScaleDownMongodbDBInstanceCpu` 方法（调用 ScaleDownDBInstanceCpu API 关闭 CPU 弹性扩容）
- [x] 1.4 为 `DescribeMongodbInstanceById` 添加重试逻辑（使用 `resource.Retry` 和 `tccommon.ReadRetryTimeout`）
- [x] 1.5 为 `ScaleUpMongodbDBInstanceCpu` 添加重试逻辑（使用 `resource.Retry` 和 `tccommon.ReadRetryTimeout`）
- [x] 1.6 为 `ScaleDownMongodbDBInstanceCpu` 添加重试逻辑（使用 `resource.Retry` 和 `tccommon.ReadRetryTimeout`）

## 2. 资源实现

- [x] 2.1 创建 `resource_tc_mongodb_instance_cpu_scale_config.go`
- [x] 2.2 实现资源 Schema 定义（1个必填字段：instance_id；1个可选字段：extra_cpu；1个输出字段：flow_id）
- [x] 2.3 实现 `resourceTencentCloudMongodbInstanceCpuScaleConfigCreate` - 设置 instance_id 为资源 ID，调用 ScaleUpDBInstanceCpu API
- [x] 2.4 实现 `resourceTencentCloudMongodbInstanceCpuScaleConfigRead` - 调用 DescribeDBInstances API 验证实例存在
- [x] 2.5 实现 `resourceTencentCloudMongodbInstanceCpuScaleConfigUpdate` - 根据 extra_cpu 变化路由到 ScaleUpDBInstanceCpu 或 ScaleDownDBInstanceCpu
- [x] 2.6 实现 `resourceTencentCloudMongodbInstanceCpuScaleConfigDelete` - 调用 ScaleDownDBInstanceCpu API 关闭弹性扩容
- [x] 2.7 添加 Import 支持（使用 instance_id 作为导入 ID，使用 `schema.ImportStatePassthrough`）
- [x] 2.8 `instance_id` 字段添加 `ForceNew: true` 标记

## 3. Provider 注册

- [x] 3.1 在 `provider.go` 的 ResourcesMap 中注册 `tencentcloud_mongodb_instance_cpu_scale` 资源
- [x] 3.2 在 `provider.md` 中添加资源链接

## 4. 测试实现

- [x] 4.1 创建 `resource_tc_mongodb_instance_cpu_scale_config_test.go`
- [x] 4.2 实现 `TestAccTencentCloudMongodbInstanceCpuScaleConfig_basic` 测试用例（使用 gomonkey mock 云 API）
- [x] 4.3 实现测试用例包含创建场景（开启弹性扩容）
- [x] 4.4 实现测试用例包含更新场景（调整 extra_cpu）
- [x] 4.5 实现测试用例包含删除场景（关闭弹性扩容）
- [x] 4.6 实现 Import 场景测试

## 5. 文档编写

- [x] 5.1 创建 `resource_tc_mongodb_instance_cpu_scale_config.md` 资源文档
- [x] 5.2 文档包含完整的使用示例（Example Usage）和 Import 说明
- [x] 5.3 在 `provider.md` 中确认资源已列出

## 6. 代码质量检查

- [x] 6.1 确保代码可编译（无需手动执行 go build）
- [x] 6.2 检查错误处理和日志记录完整性
- [x] 6.3 确保所有字段都有正确的 Description
- [x] 6.4 确保异步操作后包含 Read 轮询逻辑
- [x] 6.5 验证资源名称使用小写蛇形命名（如 `mongodb_instance_cpu_scale`）

## 注意事项

### 配置型资源特点
- 资源管理的是已存在 MongoDB 实例的 CPU 弹性扩容配置，而非独立的云资源
- `instance_id` 即为资源 ID（使用实例 ID 作为唯一标识）
- Create 设置 ID 后调用 Update 逻辑（复用 ScaleUp 调用）

### 异步操作处理
- `ScaleUpDBInstanceCpu` 和 `ScaleDownDBInstanceCpu` 为异步接口，返回 `FlowId`
- Create/Update/Delete 调用异步接口后，通过 `resource.Retry` + Read 轮询等待操作生效

### Update 路由逻辑
- 若 `extra_cpu` 存在且值变更 → 调用 ScaleUpDBInstanceCpu
- 若 `extra_cpu` 被移除（d.HasChange 检测到且新值为空/0）→ 调用 ScaleDownDBInstanceCpu
- 若 `extra_cpu` 无变化 → 无需操作

### Read 限制
- `DescribeDBInstances` API 返回的 `InstanceDetail` 中不包含弹性扩容的 `extra_cpu` 字段
- Read 方法仅验证实例存在性，不将 `extra_cpu` 从 API 读取回 state
- `extra_cpu` 配置值依赖 terraform state 维护