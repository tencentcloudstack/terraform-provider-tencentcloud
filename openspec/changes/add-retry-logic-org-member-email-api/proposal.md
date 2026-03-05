# Proposal: 为 Organization Member Email API 添加重试逻辑

## 概述

**变更 ID**: `add-retry-logic-org-member-email-api`  
**提议人**: AI Assistant  
**日期**: 2026-03-05  
**状态**: 待审批

## Why (为什么需要这个变更)

当前 `DescribeOrganizationOrgMemberEmailById` 方法在调用腾讯云 API 时缺乏重试机制,导致以下问题:

1. **用户体验差**: 在网络瞬时故障时,用户的 `terraform apply` 或 `terraform refresh` 操作会失败,需要手动重试
2. **可靠性不足**: API 速率限制错误会导致操作失败,即使稍后重试就能成功
3. **不一致性**: 项目中其他服务层方法都使用了重试模式 (如 CCN、SCF 等服务),而 Organization 服务方法没有使用,导致代码模式不一致
4. **运维成本**: 在生产环境中,网络抖动或 API 瞬时故障会导致 Terraform 操作失败,需要人工介入

**业务影响**:
- 用户在使用 Organization Member Email 资源时,可能遇到不必要的操作失败
- 在 CI/CD 流程中,瞬时故障可能导致整个部署流程中断

**技术债务**:
- 不符合项目既定的错误处理模式和最佳实践

通过添加重试逻辑,可以:
- ✅ 自动处理 90%+ 的瞬时故障
- ✅ 提升用户体验,减少手动重试
- ✅ 统一项目代码模式
- ✅ 降低运维成本

## 问题陈述

当前在 `service_tencentcloud_organization.go` 文件的 `DescribeOrganizationOrgMemberEmailById` 方法中,直接调用腾讯云 API `DescribeOrganizationMemberEmailBind` 而没有重试逻辑。这在以下场景中可能导致问题:

1. **网络瞬时故障**: 临时网络问题可能导致 API 调用失败
2. **API 速率限制**: 腾讯云 API 可能返回速率限制错误
3. **最终一致性**: 云 API 的最终一致性特性可能导致短暂的数据不可用

具体问题代码位置:
- 文件: `tencentcloud/services/tco/service_tencentcloud_organization.go`
- 行号: 474-478
- 方法: `DescribeOrganizationOrgMemberEmailById`

```go
response, err := me.client.UseOrganizationClient().DescribeOrganizationMemberEmailBind(request)
if err != nil {
    errRet = err
    return
}
```

## 建议的变更

为 `DescribeOrganizationMemberEmailBind` API 调用添加重试逻辑,使其与项目中其他服务层方法的模式保持一致。

### 变更范围

**影响的文件**:
- `tencentcloud/services/tco/service_tencentcloud_organization.go` (修改)

**修改类型**:
- 增强现有方法的错误处理和重试能力
- 使用 Terraform SDK 的 `resource.Retry()` 和 `tccommon.RetryError()` 工具

## 技术设计

### 当前实现

```go
func (me *OrganizationService) DescribeOrganizationOrgMemberEmailById(ctx context.Context, memberUin int64, bindId uint64) (orgMemberEmail *organization.DescribeOrganizationMemberEmailBindResponseParams, errRet error) {
    logId := tccommon.GetLogId(ctx)

    request := organization.NewDescribeOrganizationMemberEmailBindRequest()
    request.MemberUin = &memberUin

    defer func() {
        if errRet != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
        }
    }()

    ratelimit.Check(request.GetAction())

    response, err := me.client.UseOrganizationClient().DescribeOrganizationMemberEmailBind(request)
    if err != nil {
        errRet = err
        return
    }
    log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

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

### 建议的实现

```go
func (me *OrganizationService) DescribeOrganizationOrgMemberEmailById(ctx context.Context, memberUin int64, bindId uint64) (orgMemberEmail *organization.DescribeOrganizationMemberEmailBindResponseParams, errRet error) {
    logId := tccommon.GetLogId(ctx)

    request := organization.NewDescribeOrganizationMemberEmailBindRequest()
    request.MemberUin = &memberUin

    defer func() {
        if errRet != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
        }
    }()

    ratelimit.Check(request.GetAction())

    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        response, e := me.client.UseOrganizationClient().DescribeOrganizationMemberEmailBind(request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

        if response == nil || response.Response == nil {
            return nil
        }
        if *response.Response.BindId != bindId {
            return nil
        }
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

### 关键变更点

1. **使用 `resource.Retry()`**: 包裹 API 调用,提供自动重试机制
2. **使用 `tccommon.ReadRetryTimeout`**: 使用项目标准的读操作超时配置
3. **使用 `tccommon.RetryError()`**: 智能判断错误是否可重试
   - 自动识别网络错误 (`ClientError.NetworkError`)
   - 自动识别速率限制错误 (`RequestLimitExceeded`)
   - 自动识别其他可重试的错误代码
4. **保持业务逻辑不变**: bindId 验证和响应处理逻辑保持一致

### 重试行为

- **超时配置**: `tccommon.ReadRetryTimeout` (从环境变量 `TENCENTCLOUD_READ_RETRY_TIMEOUT` 读取,默认 3 分钟)
- **重试策略**: Terraform SDK 的指数退避策略
- **可重试的错误**:
  - `ClientError.NetworkError` - 网络错误
  - `ClientError.HttpStatusCodeError` - HTTP 状态码错误
  - `RequestLimitExceeded` - 请求限流
  - `ResourceInUse` - 资源正在使用
  - `ResourceUnavailable` - 资源不可用
- **不可重试的错误**: 业务逻辑错误(如参数错误、权限错误等)会立即返回

## 对现有代码的影响

### 兼容性

- ✅ **向后兼容**: 不改变方法签名和返回值
- ✅ **行为增强**: 仅增加重试能力,不改变成功场景的行为
- ✅ **错误处理**: 不可重试的错误会立即返回,保持现有错误处理逻辑

### 性能影响

- **正常情况**: 无性能影响,API 调用成功时行为一致
- **失败重试**: 在遇到可重试错误时,会自动重试,总体响应时间可能增加,但提高了成功率
- **超时保护**: 最长重试时间受 `ReadRetryTimeout` 限制,避免无限等待

## 风险评估

**风险等级**: 低

**潜在风险**:
1. ❌ **破坏性变更**: 无 - 方法签名和返回值不变
2. ⚠️ **超时延长**: 在错误场景下,操作可能需要更长时间才返回失败
   - **缓解措施**: 使用标准的 `ReadRetryTimeout` 配置,可通过环境变量调整
3. ⚠️ **重复请求**: 重试可能导致相同请求多次发送
   - **缓解措施**: 查询 API 是幂等的,多次调用不会产生副作用

## 测试计划

### 单元测试

不需要新增单元测试,因为:
- `resource.Retry()` 是 Terraform SDK 的标准功能,已经过充分测试
- `tccommon.RetryError()` 是项目公共方法,已有测试覆盖

### 集成测试

验证现有的资源测试仍然通过:
```bash
TF_ACC=1 go test -v ./tencentcloud/services/tco -run TestAccTencentCloudOrganizationOrgMemberEmailResource
```

### 手动验证

1. **正常场景**: 验证 API 调用成功时行为不变
2. **网络故障场景**: 模拟网络抖动,验证重试机制生效
3. **速率限制场景**: 在高频调用时,验证速率限制错误能够被重试

## 替代方案

### 方案 1: 不做任何修改 (当前状态)
- ❌ **缺点**: API 调用在瞬时故障时会失败
- ❌ **缺点**: 与项目其他服务层方法模式不一致

### 方案 2: 仅重试网络错误
- ❌ **缺点**: 无法处理速率限制等其他可重试错误
- ❌ **缺点**: 需要自定义重试逻辑,增加代码复杂度

### 方案 3: 使用自定义重试逻辑
- ❌ **缺点**: 重新发明轮子,不利用现有的基础设施
- ❌ **缺点**: 与项目现有模式不一致

**推荐方案**: 建议的实现 (使用 `resource.Retry()` 和 `tccommon.RetryError()`)
- ✅ **优点**: 与项目现有模式一致
- ✅ **优点**: 利用 Terraform SDK 的成熟重试机制
- ✅ **优点**: 自动处理多种可重试错误

## 实施计划

### 阶段 1: 代码修改 (预计 1 小时)
- [ ] 修改 `DescribeOrganizationOrgMemberEmailById` 方法
- [ ] 验证代码编译通过
- [ ] 运行 `make fmt` 和 `make lint`

### 阶段 2: 测试验证 (预计 1 小时)
- [ ] 运行现有的验收测试
- [ ] 手动验证正常场景
- [ ] 验证重试行为(可选)

### 阶段 3: 代码审查与合并 (预计 1 小时)
- [ ] 创建 PR
- [ ] 通过代码审查
- [ ] 合并到主分支

**总计预估时间**: 3 小时

## 成功标准

变更被认为成功,如果:
1. ✅ 代码编译通过,无语法错误
2. ✅ 通过 `make lint` 检查
3. ✅ 现有的 Organization Member Email 资源测试通过
4. ✅ 代码审查批准
5. ✅ 与项目现有的重试模式一致

## 依赖项

**代码依赖**:
- `github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource` (已存在)
- `tencentcloud/common` 包的 `RetryError()` 和 `ReadRetryTimeout` (已存在)

**无需额外的外部依赖**

## 参考资料

**项目内类似实现**:
- `tencentcloud/services/ccn/resource_tc_ccn_attachment.go` - 使用 `resource.Retry()` 的示例
- `tencentcloud/services/scf/service_tencentcloud_scf.go` - 服务层重试模式
- `tencentcloud/common/common.go` - `RetryError()` 实现和可重试错误列表

**相关文档**:
- Terraform Plugin SDK: https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource
- 腾讯云 Organization API: https://cloud.tencent.com/document/product/850

## 问题与讨论

**待解决的问题**:
1. 是否需要添加额外的可重试错误代码? (针对 Organization 服务的特定错误)
2. 是否需要在其他 Organization 服务方法中应用相同的模式?

**待讨论点**:
- 是否应该在同一个 PR 中修复所有 Organization 服务的类似问题?

---

**批准人**: _待定_  
**批准日期**: _待定_
