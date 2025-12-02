# Change: 添加 MongoDB 实例 SSL 配置资源

## Why

目前 Terraform Provider 已经支持 MongoDB 实例、参数、透明数据加密等配置管理，但缺少 SSL 访问配置功能。用户无法通过 Terraform 配置 MongoDB 实例的 SSL 访问状态，这限制了安全配置的自动化能力。

腾讯云提供了两个 API 接口来管理 MongoDB SSL 配置：
- `DescribeInstanceSSL` - 查询实例 SSL 开启状态
- `InstanceEnableSSL` - 设置实例 SSL 状态（开启/关闭）

通过实现 `tencentcloud_mongodb_instance_ssl` 资源，用户可以：
- 以声明式方式管理 MongoDB 实例的 SSL 访问配置
- 自动开启或关闭 SSL 加密传输
- 获取 SSL 证书下载链接和过期时间
- 实现安全合规要求的自动化配置
- 提高数据传输安全性

## What Changes

新增 Terraform 配置型资源 `tencentcloud_mongodb_instance_ssl`，支持完整的 CRUD 操作：

### 新增文件
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance_ssl.go` - 资源实现
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance_ssl_test.go` - 验收测试
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance_ssl.md` - 资源文档
- `website/docs/r/mongodb_instance_ssl.html.markdown` - 网站文档

### 修改文件
- `tencentcloud/provider.go` - 注册新资源
- `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go` - 添加 SSL 相关服务方法

### 资源 Schema
```hcl
resource "tencentcloud_mongodb_instance_ssl" "example" {
  instance_id = "cmgo-xxxxxxxx"
  enable      = true
}
```

### 字段说明
- `instance_id` (必填, ForceNew) - MongoDB 实例 ID
- `enable` (必填) - 是否开启 SSL，`true` 为开启，`false` 为关闭

### 输出属性
- `status` (只读) - SSL 开启状态，0 表示关闭，1 表示开启
- `expired_time` (只读) - 证书过期时间（格式：2023-05-01 12:00:00）
- `cert_url` (只读) - 证书下载链接（仅开启 SSL 时有值）

### 资源 ID 格式
使用实例 ID 作为资源 ID：`{instanceId}`

例如：`cmgo-p8vnipr5`

## Impact

### 受影响的规范
- 新增规范：`mongodb-ssl-config` - MongoDB SSL 配置管理

### 受影响的代码
- `tencentcloud/services/mongodb/` - 新增 SSL 配置资源实现
- `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go` - 新增服务方法
- `tencentcloud/provider.go` - 资源注册

### 向后兼容性
- ✅ 完全向后兼容，不影响现有资源
- ✅ 新增资源，不修改现有 API
- ✅ 符合配置型资源模式（Config Resource）

### 依赖关系
- 依赖 `tencentcloud_mongodb_instance` 或其他 MongoDB 实例资源 - 需要已存在的 MongoDB 实例

### 测试影响
- 需要验收测试环境中的 MongoDB 实例
- SSL 配置可能需要一定时间生效

### 类似资源参考
本资源参考以下配置型资源的实现模式：
- `tencentcloud_mongodb_instance_params` - MongoDB 参数配置
- `tencentcloud_mongodb_instance_transparent_data_encryption` - MongoDB 透明数据加密
- `tencentcloud_tdmq_rabbitmq_user_permission` - TDMQ RabbitMQ 权限配置

### 资源特性
- **配置型资源**：管理实例的某个配置项，而非独立的云资源
- **使用实例 ID 作为资源 ID**：遵循配置型资源的 ID 模式
- **无 Delete 真实操作**：删除资源时关闭 SSL（或仅从状态移除）
- **支持 Import**：可导入已存在的 SSL 配置
