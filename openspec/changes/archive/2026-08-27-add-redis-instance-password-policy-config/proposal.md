## Why

Redis 实例密码复杂度策略目前无法通过 Terraform 管理，用户只能登录腾讯云控制台手动配置。这限制了安全合规自动化能力——用户无法以基础设施即代码方式统一管理 Redis 实例的密码策略（启用开关、字母/数字/特殊字符最小数量、密码最小总长度），难以在多实例、多环境下保持一致的安全基线。

## What Changes

- 新增 Terraform 配置型资源 `tencentcloud_redis_instance_password_policy_config`，用于管理 Redis 实例的密码复杂度策略配置。
- 资源使用 `instance_id` 作为资源 ID（单例配置，资源存在即配置存在）。
- Create = SetId + Update（遵循 RESOURCE_KIND_CONFIG 模式）；Read 调用 `DescribeInstancePasswordPolicy`；Update 调用 `ModifyInstancePasswordPolicy`；Delete 为 no-op（配置随实例存在，无独立删除接口）。
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源。
- 新增资源文档 `resource_tc_redis_instance_password_policy_config.md`。
- 新增单元测试 `resource_tc_redis_instance_password_policy_config_test.go`（使用 gomonkey mock 云 API）。

## Capabilities

### New Capabilities
- `redis-instance-password-policy-config`: 管理 Redis 实例密码复杂度策略配置（启用状态、字母/数字/特殊字符最小数量、密码最小总长度）的读取与更新。

### Modified Capabilities

无。

## Impact

- **新增文件**:
  - `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config.go`
  - `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config_test.go`
  - `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config.md`
- **修改文件**:
  - `tencentcloud/provider.go` — 注册 `tencentcloud_redis_instance_password_policy_config`
  - `tencentcloud/provider.md` — 添加资源声明
- **云 API 依赖**:
  - `DescribeInstancePasswordPolicy`（redis/v20180412）— 读取密码策略
  - `ModifyInstancePasswordPolicy`（redis/v20180412）— 修改密码策略
- **向后兼容性**: 完全向后兼容，新增资源不影响现有功能。
- **异步行为**: 两个接口均为同步接口，`ModifyInstancePasswordPolicy` 不返回 TaskId，无需异步轮询。
