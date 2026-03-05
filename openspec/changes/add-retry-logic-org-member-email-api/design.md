# 技术设计: Organization Member Email API 重试逻辑

**变更 ID**: `add-retry-logic-org-member-email-api`  
**设计人**: AI Assistant  
**日期**: 2026-03-05

---

## 设计目标

为 `DescribeOrganizationOrgMemberEmailById` 服务层方法添加重试逻辑,使其能够:
1. 自动处理瞬时网络故障
2. 自动处理 API 速率限制
3. 提高资源读取操作的可靠性
4. 与项目现有的重试模式保持一致

---

## 架构概述

### 调用链路

```
Terraform Resource (Read 操作)
    ↓
resourceTencentCloudOrganizationOrgMemberEmailRead()
    ↓
service.DescribeOrganizationOrgMemberEmailById()  ← 在此添加重试
    ↓
client.UseOrganizationClient().DescribeOrganizationMemberEmailBind()
    ↓
腾讯云 Organization API
```

### 重试层级

```
┌─────────────────────────────────────────────┐
│  Resource Layer (资源层)                     │
│  - 调用服务层方法                            │
│  - 处理 Terraform 状态                       │
└─────────────────┬───────────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────────┐
│  Service Layer (服务层) ← 重试逻辑在此层    │
│  - DescribeOrganizationOrgMemberEmailById() │
│  - 包含 resource.Retry()                    │
└─────────────────┬───────────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────────┐
│  SDK Layer (SDK 层)                         │
│  - 腾讯云 SDK 客户端                         │
│  - HTTP 请求发送                            │
└─────────────────────────────────────────────┘
```

**设计原则**: 重试逻辑放在服务层,而不是资源层或 SDK 层:
- ✅ **服务层**: 提供可复用的重试能力,所有使用该服务方法的地方都受益
- ❌ **资源层**: 重复代码,每个资源都需要单独添加重试
- ❌ **SDK 层**: SDK 是第三方库,不应修改

---

## 详细设计

### 重试机制

#### 使用 Terraform SDK 的 resource.Retry()

```go
err := resource.Retry(timeout, func() *resource.RetryError {
    // API 调用
    response, e := client.SomeAPICall(request)
    
    if e != nil {
        // 智能判断是否重试
        return tccommon.RetryError(e)
    }
    
    // 业务逻辑处理
    return nil  // 成功,停止重试
})
```

**关键组件**:

1. **`resource.Retry()`**: Terraform Plugin SDK 提供的重试工具
   - 自动实现指数退避策略
   - 支持超时配置
   - 支持可重试和不可重试错误区分

2. **`tccommon.RetryError()`**: 项目自定义的错误判断函数
   - 根据错误类型判断是否可重试
   - 识别腾讯云 SDK 错误码
   - 返回 `resource.RetryableError()` 或 `resource.NonRetryableError()`

3. **`tccommon.ReadRetryTimeout`**: 项目标准超时配置
   - 默认值: 3 分钟 (180 秒)
   - 可通过环境变量 `TENCENTCLOUD_READ_RETRY_TIMEOUT` 配置
   - 适用于所有读取操作

#### 重试决策流程

```
API 调用
    ↓
是否返回错误?
    │
    ├─ 否 → 成功,返回结果
    │
    └─ 是 → tccommon.RetryError() 判断错误类型
              │
              ├─ ClientError.NetworkError → 可重试
              ├─ ClientError.HttpStatusCodeError → 可重试
              ├─ RequestLimitExceeded → 可重试
              ├─ ResourceInUse → 可重试
              ├─ ResourceUnavailable → 可重试
              ├─ ResourceBusy → 可重试 (特定服务)
              ├─ 其他预定义错误码 → 可重试
              │
              └─ 其他错误 (参数错误、权限错误等) → 不可重试
                      ↓
                 立即返回错误
```

### 可重试错误列表

在 `tencentcloud/common/common.go` 中定义:

```go
var retryableErrorCode = []string{
    // 客户端错误
    "ClientError.NetworkError",
    "ClientError.HttpStatusCodeError",
    
    // 通用错误
    "RequestLimitExceeded",
    "ResourceInUse",
    "ResourceUnavailable",
    
    // 特定服务错误
    "ResourceBusy",  // CBS
    "InvalidParameter.ActionInProgress",  // TEO
    "OperationDenied.InstanceStatusLimitError",  // PostgreSQL
    "UnsupportedOperation.UnsupportedDeleteService",  // APIGW
    "FailedOperation.ListenerHasTask",  // GAAP
}
```

**扩展机制**: `tccommon.RetryError()` 支持传入额外的可重试错误码:
```go
return tccommon.RetryError(e, "FailedOperation.SpecificError")
```

---

## 实现细节

### 代码变更对比

#### 当前实现 (无重试)

```go
func (me *OrganizationService) DescribeOrganizationOrgMemberEmailById(
    ctx context.Context, 
    memberUin int64, 
    bindId uint64,
) (orgMemberEmail *organization.DescribeOrganizationMemberEmailBindResponseParams, errRet error) {
    logId := tccommon.GetLogId(ctx)
    
    request := organization.NewDescribeOrganizationMemberEmailBindRequest()
    request.MemberUin = &memberUin

    defer func() {
        if errRet != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", 
                logId, request.GetAction(), request.ToJsonString(), errRet.Error())
        }
    }()

    ratelimit.Check(request.GetAction())

    // 直接调用 API,无重试
    response, err := me.client.UseOrganizationClient().DescribeOrganizationMemberEmailBind(request)
    if err != nil {
        errRet = err
        return
    }
    
    log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", 
        logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

    if response == nil || response.Response == nil {
        return
    }
    if *response.Response.BindId != bindId {
        return
    }
    
    orgMemberEmail = response.Response
    return
}
```

#### 新实现 (带重试)

```go
func (me *OrganizationService) DescribeOrganizationOrgMemberEmailById(
    ctx context.Context, 
    memberUin int64, 
    bindId uint64,
) (orgMemberEmail *organization.DescribeOrganizationMemberEmailBindResponseParams, errRet error) {
    logId := tccommon.GetLogId(ctx)
    
    request := organization.NewDescribeOrganizationMemberEmailBindRequest()
    request.MemberUin = &memberUin

    defer func() {
        if errRet != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", 
                logId, request.GetAction(), request.ToJsonString(), errRet.Error())
        }
    }()

    ratelimit.Check(request.GetAction())

    // 使用 resource.Retry() 包裹 API 调用
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        response, e := me.client.UseOrganizationClient().DescribeOrganizationMemberEmailBind(request)
        if e != nil {
            // 智能判断是否重试
            return tccommon.RetryError(e)
        }
        
        log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", 
            logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

        // 验证响应
        if response == nil || response.Response == nil {
            return nil
        }
        if *response.Response.BindId != bindId {
            return nil
        }
        
        // 保存结果
        orgMemberEmail = response.Response
        return nil
    })
    
    if err != nil {
        errRet = err
        return
    }

    return
}
```

### 关键变更点分析

| 方面 | 当前实现 | 新实现 | 说明 |
|------|---------|--------|------|
| **API 调用方式** | 直接调用 | `resource.Retry()` 包裹 | 添加重试能力 |
| **错误处理** | `if err != nil { return }` | `tccommon.RetryError(e)` | 智能判断是否重试 |
| **超时控制** | 无超时保护 | `tccommon.ReadRetryTimeout` | 最长 3 分钟 |
| **日志位置** | API 调用后 | 重试函数内部 | 每次尝试都记录 |
| **响应验证** | API 调用后 | 重试函数内部 | 确保验证逻辑也在重试范围内 |
| **变量作用域** | 函数级 | 闭包捕获 `orgMemberEmail` | 通过闭包传递结果 |

---

## 性能分析

### 时间复杂度

**正常场景 (API 调用成功)**:
- 当前实现: O(1) - 单次 API 调用
- 新实现: O(1) - 单次 API 调用,立即成功
- **性能影响**: 无

**失败场景 (可重试错误)**:
- 当前实现: O(1) - 立即失败返回
- 新实现: O(n) - 最多重试 n 次,直到成功或超时
  - n 由指数退避策略和超时时间决定
  - 典型值: 3-5 次重试
- **性能影响**: 增加响应时间,但提高成功率

### 指数退避策略

Terraform SDK 的默认退避策略:
```
第 1 次重试: 立即
第 2 次重试: ~1 秒后
第 3 次重试: ~2 秒后
第 4 次重试: ~4 秒后
第 5 次重试: ~8 秒后
...
最大间隔: 30 秒
总超时: tccommon.ReadRetryTimeout (3 分钟)
```

**示例时间轴** (速率限制场景):
```
0s    : 第 1 次调用 → RequestLimitExceeded
0s    : 第 2 次调用 → RequestLimitExceeded
1s    : 第 3 次调用 → RequestLimitExceeded
3s    : 第 4 次调用 → 成功
总耗时: ~3 秒
```

---

## 错误处理策略

### 错误分类

#### 1. 可重试的瞬时错误
**特征**: 重试后可能成功
**示例**:
- `ClientError.NetworkError` - 网络抖动
- `RequestLimitExceeded` - 速率限制,等待后可恢复
- `ResourceUnavailable` - 资源暂时不可用

**处理**: 自动重试,直到成功或超时

#### 2. 不可重试的业务错误
**特征**: 重试无法解决问题
**示例**:
- `InvalidParameter` - 参数错误
- `UnauthorizedOperation` - 权限不足
- `ResourceNotFound` - 资源不存在

**处理**: 立即返回错误,不重试

#### 3. 空响应或数据不匹配
**特征**: API 调用成功,但响应不符合预期
**示例**:
- `response.Response == nil` - 响应为空
- `*response.Response.BindId != bindId` - BindId 不匹配

**处理**: 返回 `nil` 错误,停止重试,由上层决定如何处理

### 错误传播

```
API 错误
    ↓
tccommon.RetryError() 判断
    ↓
┌───────────────────┐        ┌────────────────────┐
│ 可重试错误        │        │ 不可重试错误       │
│ RetryableError    │        │ NonRetryableError  │
└─────┬─────────────┘        └──────┬─────────────┘
      │                              │
      ↓                              ↓
  等待 + 重试                    立即返回
      │                              │
      ↓                              │
  达到超时或成功                     │
      │                              │
      └──────────┬───────────────────┘
                 ↓
           返回给资源层
                 ↓
           Terraform 状态更新
```

---

## 并发与线程安全

### 并发场景

**场景**: 多个 Terraform 资源同时读取 Organization Member Email
```
terraform apply (并发 = 10)
    ↓
并发创建/读取 10 个 org_member_email 资源
    ↓
同时调用 DescribeOrganizationOrgMemberEmailById()
    ↓
可能触发 API 速率限制
    ↓
重试机制自动处理
```

### 线程安全

**分析**:
1. **`resource.Retry()`**: 线程安全,每个调用独立
2. **`orgMemberEmail` 变量**: 通过闭包捕获,每个调用独立
3. **腾讯云客户端**: SDK 客户端线程安全
4. **速率限制**: `ratelimit.Check()` 在重试前调用,协调并发请求

**结论**: 设计是线程安全的,支持并发调用

---

## 监控与可观测性

### 日志输出

#### 成功场景
```
[DEBUG] [request_id] api[DescribeOrganizationMemberEmailBind] success, 
request body [...], response body [...]
```

#### 重试场景
```
[CRITAL] Retryable defined error: RequestLimitExceeded
[DEBUG] [request_id] api[DescribeOrganizationMemberEmailBind] success (重试后成功)
```

#### 失败场景
```
[CRITAL] [request_id] api[DescribeOrganizationMemberEmailBind] fail, 
request body [...], reason[InvalidParameter]
```

### 调试方法

**查看重试日志**:
```bash
TF_LOG=DEBUG terraform apply 2>&1 | grep -A 5 "Retryable"
```

**查看 API 调用详情**:
```bash
TF_LOG=DEBUG terraform apply 2>&1 | grep "DescribeOrganizationMemberEmailBind"
```

---

## 测试策略

### 单元测试
**不需要新增单元测试**,原因:
- `resource.Retry()` 是 Terraform SDK 标准功能,已测试
- `tccommon.RetryError()` 是公共函数,已有测试覆盖

### 集成测试
**运行现有的验收测试**:
```bash
TF_ACC=1 go test -v ./tencentcloud/services/tco \
    -run TestAccTencentCloudOrganizationOrgMemberEmailResource \
    -timeout 120m
```

### 手动测试场景

#### 场景 1: 正常调用
**目的**: 验证正常场景下行为不变
**步骤**:
1. 创建 Organization Member Email 资源
2. 观察日志,确认单次 API 调用成功
3. 验证资源状态正确

#### 场景 2: 网络抖动模拟
**目的**: 验证重试机制处理网络错误
**步骤**:
1. 在本地网络引入延迟或丢包 (tcpdump/iptables)
2. 执行 terraform refresh
3. 观察日志中的重试行为
4. 验证最终成功

#### 场景 3: 速率限制模拟
**目的**: 验证重试机制处理速率限制
**步骤**:
1. 并发创建多个资源 (触发速率限制)
2. 观察日志中的 `RequestLimitExceeded` 错误和重试
3. 验证所有资源最终创建成功

---

## 向后兼容性

### API 兼容性
✅ **完全兼容**
- 方法签名不变
- 返回值类型不变
- 行为增强 (仅增加重试,不改变成功场景逻辑)

### 配置兼容性
✅ **完全兼容**
- 现有 Terraform 配置无需修改
- 环境变量 `TENCENTCLOUD_READ_RETRY_TIMEOUT` 可选配置

### 状态文件兼容性
✅ **完全兼容**
- 资源 ID 格式不变
- 状态数据结构不变

---

## 安全性考虑

### 敏感信息
**问题**: 重试时可能多次记录请求/响应日志,是否泄露敏感信息?

**分析**:
- `DescribeOrganizationMemberEmailBind` 是查询 API
- 响应包含邮箱、手机号等信息
- 日志级别: DEBUG (生产环境默认不开启)

**缓解措施**:
- 日志仅在 `TF_LOG=DEBUG` 时输出
- 敏感字段应在 Terraform Schema 中标记 `Sensitive: true`
- 建议不在生产环境开启 DEBUG 日志

### 重放攻击
**问题**: 重试是否可能导致重放攻击?

**分析**:
- `DescribeOrganizationMemberEmailBind` 是幂等的查询操作
- 多次调用不会产生副作用
- 请求签名由 SDK 自动生成,包含时间戳

**结论**: 无重放攻击风险

---

## 替代设计方案

### 方案 A: 在资源层添加重试 (未采用)
```go
func resourceTencentCloudOrganizationOrgMemberEmailRead(d *schema.ResourceData, meta interface{}) error {
    // 在资源层使用 resource.Retry()
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        orgMemberEmail, e := service.DescribeOrganizationOrgMemberEmailById(...)
        if e != nil {
            return tccommon.RetryError(e)
        }
        // 设置状态
        return nil
    })
}
```

**缺点**:
- ❌ 需要在每个使用该服务方法的地方重复添加重试
- ❌ 代码重复,维护成本高
- ❌ 与项目现有模式不一致

### 方案 B: 在 SDK 层添加重试 (未采用)
修改腾讯云 SDK 客户端配置,启用自动重试

**缺点**:
- ❌ SDK 是第三方依赖,不应修改
- ❌ SDK 的重试策略可能不符合 Terraform 场景
- ❌ 难以与 Terraform 的状态管理集成

### 方案 C: 自定义重试逻辑 (未采用)
```go
for i := 0; i < maxRetries; i++ {
    response, err := client.DescribeOrganizationMemberEmailBind(request)
    if err == nil {
        break
    }
    if !isRetryable(err) {
        return err
    }
    time.Sleep(backoff(i))
}
```

**缺点**:
- ❌ 重新发明轮子,维护成本高
- ❌ 需要自己实现退避策略、超时控制等
- ❌ 与项目现有工具不一致

---

## 未来扩展

### 短期 (1-3 个月)
- [ ] 审查 Organization 服务的其他方法,应用相同的重试模式
- [ ] 收集生产环境的重试指标 (成功率、重试次数等)

### 中期 (3-6 个月)
- [ ] 基于实际数据优化超时配置
- [ ] 考虑为特定 API 添加自定义重试逻辑

### 长期 (6-12 个月)
- [ ] 跨服务统一重试模式审计
- [ ] 考虑引入更高级的重试策略 (如 jitter)

---

## 参考资料

**项目内参考**:
- `tencentcloud/common/common.go` - `RetryError()` 实现
- `tencentcloud/services/scf/service_tencentcloud_scf.go` - 服务层重试示例
- `tencentcloud/services/ccn/resource_tc_ccn_attachment.go` - 资源层重试示例

**外部文档**:
- [Terraform Plugin SDK - resource.Retry()](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource#Retry)
- [腾讯云 API 错误码](https://cloud.tencent.com/document/api/error-center)
- [Exponential Backoff](https://en.wikipedia.org/wiki/Exponential_backoff)

---

**审查人**: _待定_  
**审查日期**: _待定_
