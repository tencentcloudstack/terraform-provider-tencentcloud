## Why

腾讯云 MongoDB 提供了 CPU 弹性扩容能力，允许用户在业务高峰期手动开启 CPU 弹性扩容并设置额外 CPU 核数，在低峰期手动关闭弹性扩容回缩 CPU。目前 Terraform Provider 尚未支持通过声明式方式管理 MongoDB 实例的 CPU 弹性扩容配置，用户无法通过 Terraform 自动化管理这一配置。通过实现 `tencentcloud_mongodb_instance_cpu_scale` 配置型资源，用户可以通过 Terraform 以 IaC 方式管理 MongoDB 实例的 CPU 弹性扩容开关及额外 CPU 核数。

## What Changes

新增 Terraform 配置型资源（RESOURCE_KIND_CONFIG）`tencentcloud_mongodb_instance_cpu_scale`，支持读（Read）和更新（Update）操作：

- 新增资源文件 `resource_tc_mongodb_instance_cpu_scale_config.go`
- 新增服务层方法：调用 `ScaleUpDBInstanceCpu`、`ScaleDownDBInstanceCpu`、`DescribeDBInstances` 三个云 API
- 在 `provider.go` 中注册新资源
- 新增资源文档和测试文件
- 支持 Import 已有 CPU 弹性扩容配置

### 资源 Schema
```hcl
resource "tencentcloud_mongodb_instance_cpu_scale" "example" {
  instance_id = "cmgo-xxxxxxxx"
  extra_cpu   = 2
}
```

### 核心字段
- `instance_id` (必填, ForceNew) - MongoDB 实例 ID
- `extra_cpu` (可选) - 弹性扩容的额外 CPU 核数，设置该值表示开启弹性扩容，不设置表示关闭弹性扩容

### 输出属性
- `flow_id` (只读) - 操作流程 ID（来自 ScaleUp/ScaleDown API 返回）

## Capabilities

### New Capabilities
- `mongodb-cpu-scale-config`: MongoDB 实例 CPU 弹性扩容配置管理，支持通过 ScaleUpDBInstanceCpu 开启弹性扩容、ScaleDownDBInstanceCpu 关闭弹性扩容、DescribeDBInstances 读取实例详情

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## Impact

### 受影响的代码
- `tencentcloud/services/mongodb/` - 新增 CPU 弹性扩容配置资源实现及服务层方法
- `tencentcloud/provider.go` - 资源注册

### 向后兼容性
- ✅ 完全向后兼容，不影响现有资源
- ✅ 新增资源，不修改现有 API

### 依赖关系
- 依赖已存在的 MongoDB 实例（由 `tencentcloud_mongodb_instance` 或其他方式创建）
- 使用的三个云 API（ScaleUpDBInstanceCpu、ScaleDownDBInstanceCpu、DescribeDBInstances）均已在 vendor SDK 中可用

### API 异步处理
- `ScaleUpDBInstanceCpu` 和 `ScaleDownDBInstanceCpu` 为异步接口，返回 FlowId，更新后需要调用 Read 接口轮询直到操作生效