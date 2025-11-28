# Change: 添加 TDMQ RabbitMQ 权限管理资源

## Why

目前 Terraform Provider 已经支持 TDMQ RabbitMQ 的实例、用户和虚拟主机管理，但缺少权限管理功能。用户无法通过 Terraform 配置 RabbitMQ 用户在特定 VirtualHost 下的访问权限（配置、写入、读取权限），这限制了完整的基础设施即代码能力。

腾讯云提供了三个 API 接口来管理 RabbitMQ 权限：
- `DescribeRabbitMQPermission` - 查询权限列表
- `ModifyRabbitMQPermission` - 修改权限
- `DeleteRabbitMQPermission` - 删除权限

通过实现 `tencentcloud_tdmq_rabbitmq_user_permission` 资源，用户可以：
- 以声明式方式管理用户权限
- 控制用户在不同 VirtualHost 下的细粒度访问权限
- 遵循最小权限原则，提高安全性
- 实现完整的 RabbitMQ 基础设施自动化

## What Changes

新增 Terraform 资源 `tencentcloud_tdmq_rabbitmq_user_permission`，支持完整的 CRUD 操作：

### 新增文件
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_user_permission.go` - 资源实现
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_user_permission_test.go` - 验收测试
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_user_permission.md` - 资源文档
- `website/docs/r/tdmq_rabbitmq_user_permission.html.markdown` - 网站文档

### 修改文件
- `tencentcloud/provider.go` - 注册新资源
- `tencentcloud/services/tdmq/service_tencentcloud_tdmq.go` - 添加权限相关服务方法

### 资源 Schema
```hcl
resource "tencentcloud_tdmq_rabbitmq_user_permission" "example" {
  instance_id    = "amqp-xxxxxxxx"
  user           = "admin"
  virtual_host   = "testvhost"
  config_regexp  = ".*"
  write_regexp   = ".*"
  read_regexp    = ".*"
}
```

### 字段说明
- `instance_id` (必填) - RabbitMQ 实例 ID
- `user` (必填) - 用户名
- `virtual_host` (必填) - VirtualHost 名称
- `config_regexp` (必填) - 配置权限正则表达式，控制可声明的资源
- `write_regexp` (必填) - 写入权限正则表达式，控制可写入的资源
- `read_regexp` (必填) - 读取权限正则表达式，控制可读取的资源

### 资源 ID 格式
使用三段式复合 ID：`{instanceId}#{user}#{virtualHost}`

例如：`amqp-2ppxx4rq#admin#testvhost`

## Impact

### 受影响的规范
- 新增规范：`tdmq-rabbitmq-permission` - RabbitMQ 权限管理

### 受影响的代码
- `tencentcloud/services/trabbit/` - 新增权限资源实现
- `tencentcloud/services/tdmq/service_tencentcloud_tdmq.go` - 新增服务方法
- `tencentcloud/provider.go` - 资源注册

### 向后兼容性
- ✅ 完全向后兼容，不影响现有资源
- ✅ 新增资源，不修改现有 API

### 依赖关系
- 依赖 `tencentcloud_tdmq_rabbitmq_vip_instance` - 需要已存在的 RabbitMQ 实例
- 依赖 `tencentcloud_tdmq_rabbitmq_user` - 需要已存在的用户
- 依赖 `tencentcloud_tdmq_rabbitmq_virtual_host` - 需要已存在的虚拟主机

### 测试影响
- 需要验收测试环境中的 RabbitMQ 实例
- 测试需要创建用户和虚拟主机作为前置条件
