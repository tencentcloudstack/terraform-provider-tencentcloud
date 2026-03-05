# Capability: API 重试与容错能力

**能力名称**: `api-retry-resilience`  
**所属领域**: 服务层 (Service Layer)  
**状态**: 增强中

---

## ADDED Requirements

### Requirement: Organization Member Email API 查询操作 MUST 支持自动重试

**优先级**: High  
**类型**: Reliability Enhancement  

**描述**:
`DescribeOrganizationOrgMemberEmailById` 服务层方法在调用腾讯云 `DescribeOrganizationMemberEmailBind` API 时,MUST 实现自动重试机制以提高可靠性。重试机制 SHALL 能够智能识别可重试的瞬时错误(如网络错误、速率限制),并在超时范围内自动重试。

#### Scenario: API 调用遇到网络瞬时故障时自动重试成功

**Given** (前置条件):
- Organization Member Email 资源已存在
- 用户执行 `terraform refresh` 或 `terraform apply` 触发资源读取
- 网络出现瞬时故障 (如丢包、连接超时)

**When** (操作):
- Terraform 调用 `DescribeOrganizationOrgMemberEmailById` 方法
- 第一次 API 调用因网络错误失败,返回 `ClientError.NetworkError`
- 重试机制识别为可重试错误
- 等待指数退避时间后自动重试
- 第二次 API 调用成功

**Then** (预期结果):
- 资源状态成功读取
- 用户无需手动重试
- 日志中记录重试行为
- Terraform 操作正常完成

**验收标准**:
```bash
# 测试命令
TF_LOG=DEBUG terraform refresh 2>&1 | grep "Retryable.*NetworkError"

# 预期输出包含
[CRITAL] Retryable defined error: ClientError.NetworkError
[DEBUG] api[DescribeOrganizationMemberEmailBind] success
```

---

#### Scenario: API 调用遇到速率限制时自动重试成功

**Given** (前置条件):
- 用户并发创建/读取多个 Organization Member Email 资源
- 腾讯云 API 返回速率限制错误

**When** (操作):
- Terraform 并发调用多个 `DescribeOrganizationOrgMemberEmailById` 方法
- 部分 API 调用因速率限制失败,返回 `RequestLimitExceeded`
- 重试机制识别为可重试错误
- 等待指数退避时间后自动重试
- 重试时速率限制已解除,调用成功

**Then** (预期结果):
- 所有资源状态最终成功读取
- 速率限制错误被自动处理
- 用户无需调整并发配置
- 日志中记录重试和成功信息

**验收标准**:
```bash
# 测试场景: 并发创建 5 个资源
terraform apply -parallelism=5

# 预期: 所有资源成功创建,可能有重试日志
TF_LOG=DEBUG terraform apply 2>&1 | grep -c "RequestLimitExceeded" # >= 0
TF_LOG=DEBUG terraform apply 2>&1 | grep -c "Apply complete!" # == 1
```

---

#### Scenario: API 调用遇到不可重试错误时立即失败

**Given** (前置条件):
- 用户提供无效的参数 (如不存在的 memberUin)
- 或用户权限不足

**When** (操作):
- Terraform 调用 `DescribeOrganizationOrgMemberEmailById` 方法
- API 返回业务错误,如 `InvalidParameter.MemberUinNotFound`
- 重试机制识别为不可重试错误

**Then** (预期结果):
- 错误立即返回,不进行重试
- 用户收到清晰的错误信息
- 不浪费时间在无意义的重试上
- Terraform 操作快速失败

**验收标准**:
```bash
# 测试: 使用无效的 member_uin
terraform apply

# 预期: 立即失败,无重试日志
# 错误信息包含 "InvalidParameter" 或 "ResourceNotFound"
# 不包含 "Retryable" 日志
```

---

#### Scenario: API 调用重试超时后返回最后的错误

**Given** (前置条件):
- 网络持续故障或 API 服务不可用
- 重试时间超过配置的超时时间 (默认 3 分钟)

**When** (操作):
- Terraform 调用 `DescribeOrganizationOrgMemberEmailById` 方法
- API 调用持续失败
- 重试机制在超时范围内不断重试
- 达到 `tccommon.ReadRetryTimeout` 超时时间

**Then** (预期结果):
- 返回最后一次的错误信息
- 用户收到明确的失败提示
- 避免无限等待
- 日志中记录所有重试尝试

**验收标准**:
```bash
# 模拟: 断开网络连接
terraform refresh

# 预期: 3 分钟后超时
# 日志包含多次 "Retryable" 重试记录
# 最终返回错误: "context deadline exceeded" 或最后的网络错误
```

---

#### Scenario: 正常 API 调用无性能影响

**Given** (前置条件):
- 网络状况良好
- API 服务正常
- 无速率限制

**When** (操作):
- Terraform 调用 `DescribeOrganizationOrgMemberEmailById` 方法
- API 第一次调用成功

**Then** (预期结果):
- 行为与添加重试前完全一致
- 无额外的延迟
- 无重试日志
- 性能无影响

**验收标准**:
```bash
# 测试: 正常读取资源
time terraform refresh

# 预期:
# - 执行时间与未添加重试前相当
# - 日志中只有一次 API 调用成功记录
# - 无 "Retryable" 日志
```

---

## 技术规范

### 实现要求

**必须 (MUST)**:
1. 使用 `resource.Retry()` 包裹 API 调用
2. 使用 `tccommon.ReadRetryTimeout` 作为超时配置
3. 使用 `tccommon.RetryError()` 判断错误是否可重试
4. 日志输出必须在重试函数内部,记录每次尝试
5. 响应验证逻辑必须在重试函数内部
6. 方法签名和返回值不得改变 (保持向后兼容)

**应该 (SHOULD)**:
1. 遵循项目现有的重试模式
2. 复用 `tencentcloud/common/common.go` 中定义的可重试错误列表
3. 在 defer 函数中记录最终错误

**可以 (MAY)**:
1. 为 Organization 服务添加特定的可重试错误码 (如果需要)

### 可重试错误清单

以下错误码被认为是可重试的:
- `ClientError.NetworkError` - 网络错误
- `ClientError.HttpStatusCodeError` - HTTP 状态码错误
- `RequestLimitExceeded` - 请求限流
- `ResourceInUse` - 资源正在使用
- `ResourceUnavailable` - 资源不可用
- `ResourceBusy` - 资源繁忙 (特定服务)
- 其他在 `tencentcloud/common/common.go` 中定义的错误

### 性能要求

1. **正常场景**: 响应时间不增加 (单次 API 调用)
2. **重试场景**: 最大超时时间为 `ReadRetryTimeout` (默认 3 分钟)
3. **退避策略**: 遵循 Terraform SDK 的指数退避算法
4. **并发支持**: 重试机制必须线程安全,支持并发调用

### 日志要求

1. **成功日志**: `[DEBUG] api[...] success, request body [...], response body [...]`
2. **重试日志**: `[CRITAL] Retryable defined error: <error_code>`
3. **失败日志**: `[CRITAL] api[...] fail, request body [...], reason[...]`

---

## 依赖关系

**依赖的能力**:
- 无 (独立能力增强)

**被依赖的能力**:
- Organization Member Email 资源读取操作

**外部依赖**:
- Terraform Plugin SDK v2 (`resource.Retry()`)
- 腾讯云 SDK Go (`organization` 包)
- 项目公共包 (`tccommon.RetryError()`, `tccommon.ReadRetryTimeout`)

---

## 测试要求

### 必需的测试

1. **集成测试**: 运行现有的 `TestAccTencentCloudOrganizationOrgMemberEmailResource` 验收测试
2. **手动验证**: 正常场景和重试场景的手动验证

### 可选的测试

1. **压力测试**: 并发场景下的速率限制处理
2. **故障注入**: 模拟网络故障验证重试行为

---

## 影响范围

**修改的文件**:
- `tencentcloud/services/tco/service_tencentcloud_organization.go` (1 个方法)

**影响的资源**:
- `tencentcloud_organization_org_member_email` (Read 操作)

**影响的用户操作**:
- `terraform refresh`
- `terraform apply` (隐式读取)
- `terraform plan` (隐式读取)

---

## 安全考虑

1. **敏感信息**: 日志可能包含邮箱、手机号,仅在 DEBUG 模式下输出
2. **重放攻击**: 查询 API 是幂等的,无重放攻击风险
3. **超时保护**: 有明确的超时时间,防止资源耗尽

---

## 合规性

**向后兼容性**: ✅ 完全兼容
- 方法签名不变
- 返回值不变
- 行为增强 (仅增加重试能力)

**破坏性变更**: ❌ 无破坏性变更

---

## 参考实现

**项目内参考**:
```go
// tencentcloud/services/scf/service_tencentcloud_scf.go
err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
    ratelimit.Check(request.GetAction())
    response, err := client.GetFunction(request)
    if err != nil {
        return tccommon.RetryError(err)
    }
    // 处理响应
    return nil
})
```

---

**审查人**: _待定_  
**审查日期**: _待定_
