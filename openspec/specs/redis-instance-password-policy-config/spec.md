# redis-instance-password-policy-config Specification

## Purpose
TBD - created by archiving change add-redis-instance-password-policy-config. Update Purpose after archive.
## Requirements
### Requirement: 资源 Schema 定义

资源 `tencentcloud_redis_instance_password_policy_config` SHALL 定义以下 schema 字段：
- `instance_id`（String, Required, ForceNew）— Redis 实例 ID
- `enabled`（Bool, Required）— 是否启用密码复杂度策略
- `min_letter_count`（Int, Optional, Computed）— 大小写字母最小字符数，取值范围 [1,16]
- `min_digit_count`（Int, Optional, Computed）— 数字字符最小字符数，取值范围 [1,16]
- `min_special_count`（Int, Optional, Computed）— 特殊字符最小字符数，取值范围 [1,16]
- `min_length`（Int, Optional, Computed）— 密码最小总长度，取值范围 [8,64]

字段 SHALL 平铺在 schema 顶层，不使用 `password_policy` 嵌套对象层。

#### Scenario: Schema 包含所有字段
- **WHEN** 定义资源 schema
- **THEN** schema 包含 `instance_id`、`enabled`、`min_letter_count`、`min_digit_count`、`min_special_count`、`min_length` 六个顶层字段

#### Scenario: instance_id 为 ForceNew
- **WHEN** 用户在 apply 后修改 `instance_id`
- **THEN** Terraform 销毁并重建资源

#### Scenario: 可选计数字段支持不设置
- **WHEN** 用户未设置 `min_letter_count`、`min_digit_count`、`min_special_count`、`min_length`
- **THEN** 这些字段在 schema 中为 Optional + Computed，由云 API 返回值填充

### Requirement: Create 操作

Create 操作 SHALL 设置资源 ID 为 `instance_id`，然后调用 Update 操作完成配置写入。

#### Scenario: 首次创建配置
- **WHEN** 用户 apply 一个新的 `tencentcloud_redis_instance_password_policy_config` 资源
- **THEN** 资源 ID 被设置为 `instance_id` 的值
- **AND** 调用 `ModifyInstancePasswordPolicy` 写入密码策略配置
- **AND** 调用 Read 刷新状态

### Requirement: Read 操作

Read 操作 SHALL 调用 `DescribeInstancePasswordPolicy` 读取当前密码策略，并设置所有字段到 state。

#### Scenario: 成功读取密码策略
- **WHEN** Read 调用 `DescribeInstancePasswordPolicy` 成功返回
- **THEN** 将 `enabled`、`min_letter_count`、`min_digit_count`、`min_special_count`、`min_length` 设置到 state
- **AND** 在设置每个字段前检查云 API 返回值是否为 nil，为 nil 则不设置

#### Scenario: 实例不存在
- **WHEN** Read 调用 `DescribeInstancePasswordPolicy` 返回实例不存在或空响应
- **THEN** 先打印 `log.Printf("[CRUD] redis_instance_password_policy_config id=%s", d.Id())` 保留现场
- **AND** 调用 `d.SetId("")` 从 state 移除资源

### Requirement: Update 操作

Update 操作 SHALL 调用 `ModifyInstancePasswordPolicy` 更新密码策略，请求中包含 `InstanceId` 和 `PasswordPolicy` 对象。

#### Scenario: 成功更新密码策略
- **WHEN** Update 调用 `ModifyInstancePasswordPolicy` 成功
- **THEN** 调用 Read 刷新状态
- **AND** 不需要异步轮询（接口为同步）

#### Scenario: 更新失败可重试
- **WHEN** Update 调用 `ModifyInstancePasswordPolicy` 返回网络错误等可重试错误
- **THEN** 使用 `tccommon.RetryError` 包装错误，由外层 retry 重试

#### Scenario: 更新失败不可重试
- **WHEN** Update 调用 `ModifyInstancePasswordPolicy` 返回实例不存在等不可重试错误
- **THEN** 返回 `NonRetryableError`

### Requirement: Delete 操作

Delete 操作 SHALL 为 no-op，不调用任何云 API。

#### Scenario: 删除资源
- **WHEN** 用户执行 terraform destroy
- **THEN** 资源从 state 移除，不调用云 API

### Requirement: Import 支持

资源 SHALL 支持 Import，使用 `instance_id` 作为导入 ID。

#### Scenario: 导入已存在的密码策略配置
- **WHEN** 用户执行 `terraform import tencentcloud_redis_instance_password_policy_config.xxx <instance_id>`
- **THEN** 资源被导入到 state，调用 Read 填充字段

### Requirement: Provider 注册

资源 SHALL 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册为 `tencentcloud_redis_instance_password_policy_config`。

#### Scenario: 资源已注册
- **WHEN** 检查 provider.go
- **THEN** ResourcesMap 中存在 `tencentcloud_redis_instance_password_policy_config` 映射到 `crs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()`

### Requirement: 单元测试

资源 SHALL 提供单元测试文件 `resource_tc_redis_instance_password_policy_config_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试。

#### Scenario: 单元测试覆盖 CRUD
- **WHEN** 运行单元测试
- **THEN** 测试覆盖 Create、Read、Update、Delete 的业务逻辑
- **AND** 使用 gomonkey mock `DescribeInstancePasswordPolicy` 和 `ModifyInstancePasswordPolicy` 接口

