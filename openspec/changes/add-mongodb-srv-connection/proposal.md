# Change: 添加 MongoDB SRV 连接 URL 配置资源

## Why

目前 Terraform Provider 已经支持 MongoDB 实例、SSL、透明数据加密等配置管理，但缺少 SRV 连接 URL 配置功能。用户无法通过 Terraform 管理 MongoDB 实例的 SRV 连接 URL，这限制了现代化连接方式的自动化能力。

腾讯云提供了完整的 API 接口来管理 MongoDB SRV 连接配置：
- `EnableSRVConnectionUrl` - 开启 SRV 连接 URL（异步任务）
- `DescribeSRVConnectionDomain` - 查询 SRV 连接域名信息
- `ModifySRVConnectionUrl` - 修改 SRV 连接 URL（异步任务）
- `DisableSRVConnectionUrl` - 关闭 SRV 连接 URL
- `DescribeAsyncRequestInfo` - 查询异步任务状态

通过实现 `tencentcloud_mongodb_instance_srv_connection` 资源，用户可以：
- 以声明式方式管理 MongoDB 实例的 SRV 连接 URL
- 自动开启或关闭 SRV 连接功能
- 配置自定义域名（可选）
- 实现现代化 MongoDB 连接方式的自动化配置
- 简化客户端连接配置管理

## What Changes

新增 Terraform 配置型资源 `tencentcloud_mongodb_instance_srv_connection`，支持完整的 CRUD 操作：

### 新增文件
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance_srv_connection.go` - 资源实现
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance_srv_connection_test.go` - 验收测试
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance_srv_connection.md` - 资源文档
- `website/docs/r/mongodb_instance_srv_connection.html.markdown` - 网站文档

### 修改文件
- `tencentcloud/provider.go` - 注册新资源
- `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go` - 添加 SRV 连接相关服务方法

### 资源 Schema
```hcl
resource "tencentcloud_mongodb_instance_srv_connection" "example" {
  instance_id = "cmgo-xxxxxxxx"
  domain      = "example.mongodb.com"  # 可选参数，Optional + Computed
}
```

### 字段说明
- `instance_id` (必填, ForceNew) - MongoDB 实例 ID
- `domain` (可选 + 计算, Optional + Computed) - 自定义域名
  - 用户可以不填，系统会使用默认域名并在 Read 时填充此字段
  - 用户可以指定自定义域名
  - 修改此字段会触发 Update 操作

### 输出属性
- `srv_url` (只读, Computed) - SRV 连接 URL
- `domain` (可选 + 计算) - 实际使用的域名（用户输入或系统默认）

### 资源 ID 格式
使用实例 ID 作为资源 ID：`{instanceId}`

例如：`cmgo-p8vnipr5`

### 操作特性
- **异步操作处理**：Create 和 Update 操作会触发异步任务，需要通过 `DescribeAsyncRequestInfo` 轮询任务状态直到成功或失败
- **超时设置**：异步任务默认超时时间为 3 倍 ReadRetryTimeout
- **幂等性**：资源创建和更新操作需要处理幂等性

## Impact

### 受影响的规范
- 新增规范：`mongodb-srv-connection` - MongoDB SRV 连接 URL 配置管理

### 受影响的代码
- `tencentcloud/services/mongodb/` - 新增 SRV 连接配置资源实现
- `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go` - 新增服务方法
- `tencentcloud/provider.go` - 资源注册

### 向后兼容性
- ✅ 完全向后兼容，不影响现有资源
- ✅ 新增资源，不修改现有 API
- ✅ 符合配置型资源模式（Config Resource）

### 依赖关系
- 依赖 `tencentcloud_mongodb_instance` 或其他 MongoDB 实例资源 - 需要已存在的 MongoDB 实例
- 依赖现有的 `DescribeAsyncRequestInfo` 服务方法处理异步任务

### 测试影响
- 需要验收测试环境中的 MongoDB 实例
- 异步任务可能需要较长时间等待完成
- 需要测试自定义域名和默认域名两种场景

### 类似资源参考
本资源参考以下配置型资源的实现模式：
- `tencentcloud_mongodb_instance_ssl` - MongoDB SSL 配置
- `tencentcloud_mongodb_instance_transparent_data_encryption` - MongoDB 透明数据加密
- `tencentcloud_mongodb_instance_backup` - MongoDB 备份配置

### 资源特性
- **配置型资源**：管理实例的某个配置项，而非独立的云资源
- **使用实例 ID 作为资源 ID**：遵循配置型资源的 ID 模式
- **异步任务处理**：Create 和 Update 操作需要等待异步任务完成
- **支持 Import**：可导入已存在的 SRV 连接配置
- **Delete 操作**：删除资源时调用 DisableSRVConnectionUrl 关闭功能
