## Context

腾讯云点播（VOD）近期上线了 AIGC API Token 能力（2026-02-02 上线），用于为 AIGC 类 API 调用提供独立的凭证管理。云端提供 3 个无状态接口：

| 云 API | 入参 | 出参 |
|---|---|---|
| `CreateAigcApiToken` | `SubAppId` (必填) | `ApiToken` (String) |
| `DescribeAigcApiTokens` | `SubAppId` (可选) | `ApiTokens` (Array of String) |
| `DeleteAigcApiToken` | `SubAppId` (必填) + `ApiToken` (必填) | 无业务字段 |

当前 Provider 内：
- `tencentcloud-sdk-go/tencentcloud/vod/v20180717` 已 vendored 对应 `*Request/*Response` 与 client 方法（已在 vendor 源码确认）
- `tencentcloud/services/vod/service_tencentcloud_vod.go` 为既有 VodService 统一服务层
- `tencentcloud/services/vod/resource_tc_vod_sub_application.go` 与 `resource_tc_vod_event_config.go` 提供 VOD 资源的成熟参考模式

约束:
- 必须使用 vendor 内已有 SDK，不引入新依赖
- 必须遵循项目既有的 VOD service 层组织与错误处理模式
- 必须保持向后兼容
- 必须有完整文档模板与验收测试
- `ApiToken` 属敏感凭证，不可明文打印至日志

## Goals / Non-Goals

**Goals:**

- 提供 `tencentcloud_vod_aigc_api_token` Terraform 资源，覆盖 Create/Read/Delete 全流程
- 合理处理云端 ~30 秒的数据同步延迟（官方文档明确说明：Create/Delete 后约 30 秒才可查询到最新数据）
- Schema 字段严格对齐 SDK：`sub_app_id`（Required, ForceNew）、`api_token`（Computed, ForceNew, Sensitive）
- 支持资源导入（`terraform import`），入参为复合 ID `sub_app_id#api_token`
- 遵循项目现有 operation/resource 命名与文件组织约定
- 提供至少 1 个基础验收测试用例（创建 + 自动销毁）

**Non-Goals:**

- 不实现 Update：SDK 无对应 Modify 接口；Token 本身不可变
- 不实现数据源 `tencentcloud_vod_aigc_api_tokens`：本次仅交付资源；若后续用户提出需求再增量交付
- 不在 Provider 端缓存 Token 值以规避查询延迟（无状态重试更简单且正确）
- 不做跨 SubAppId 的批量管理：每个资源实例对应一个 Token
- 不引入新的 SDK 依赖

## Decisions

### Decision 1: 采用标准 Resource 模式（非 operation 资源）

**选择**: 使用标准 CRUD Resource（`ResourceTencentCloudVodAigcApiToken`），而不是 `_operation` 后缀资源。

**理由**:
- 该 API 有完整的 Create + List + Delete，支持对 Token 的生命周期管理
- Token 是持久化云侧实体（非一次性操作），用户需要跟踪其存在性
- Read 可通过 `DescribeAigcApiTokens` 列表中检索 token 值来实现幂等一致性

**替代方案**:
- operation 资源：不合适——API 明确提供 Describe/Delete，有对应逆向能力
- 在 `tencentcloud_vod_sub_application` 资源中以子属性形式暴露：侵入性大，且一个 SubAppId 可以创建多个 Token，子属性无法优雅表达 1:N 关系

### Decision 2: Schema 设计

| Terraform 字段 | 类型 | Required/Optional | ForceNew | Computed | Sensitive | 对应 API 字段 |
|---|---|---|---|---|---|---|
| `sub_app_id` | Int | Required | ✅ | — | — | `SubAppId` |
| `api_token` | String | Optional(Computed) | ✅ | ✅ | ✅ | `ApiToken`（响应返回） |

**关键决策**:
- `sub_app_id`: `Required + ForceNew`。Token 绑定到具体应用，变更等于换一个资源。
- `api_token`: `Computed + ForceNew + Sensitive`。Create 时由云端返回；`Sensitive: true` 避免在 plan/apply 输出与日志中明文暴露；`ForceNew` 是为了支持 `terraform import` 时用户传入已有 Token 后，若后续修改则触发重建（实际普通 apply 流程下用户不应设置该字段）。
- 不增加 `description` / `name` 字段：API 不支持。
- 不加自定义 Timeouts：Create/Delete 单次 API 调用本身很快，30 秒同步延迟由 Read 侧的 `helper.Retry` 处理。

**替代方案**:
- 将 `api_token` 设为 `Computed`（无 `Optional`）：则无法支持 `terraform import`（导入时需要指定 token 值）
- 将 `sub_app_id` 设为 `Optional`：偏离 API 语义且 state 键值不稳定

### Decision 3: 复合资源 ID

**选择**: `d.SetId(fmt.Sprintf("%d%s%s", subAppId, FILED_SP, apiToken))`，分隔符使用项目既有常量 `tccommon.FILED_SP`（`#`）。

**理由**:
- Token 本身在同一 SubAppId 下唯一，但跨 SubAppId 可能重复
- 与项目其他资源的复合 ID 模式保持一致（如 `tencentcloud_dcdb_account`: `instanceId#userName#host`）
- Import 时可直接解析复合 ID 还原两字段

**替代方案**:
- 仅用 `apiToken` 作 ID：丢失 subAppId 上下文，Read 时需额外字段
- JSON 序列化 ID：不符合项目风格

### Decision 4: 处理 30 秒数据同步延迟

**选择**: 
- **Create 后**: 先 `d.SetId()`，再调用 `resourceTencentCloudVodAigcApiTokenRead` 之前用 `helper.Retry(ctx, readRetry, func)` 轮询 `DescribeAigcApiTokens` 直到返回列表中包含刚创建的 token（超时 `tccommon.ReadRetryTimeout`，即 3 分钟）
- **Delete 后**: 同理轮询直到列表中不再包含该 token
- **Read**: 若列表中不存在该 token，设 `d.SetId("")` 交由 Terraform 感知为"已消失"

**理由**:
- 官方文档明确提示 ~30s 同步延迟，不处理会导致 Create 后立即 Read 返回空 → 被 Terraform 判定为资源消失 → 下次 apply 死循环
- 所有轮询动作在 Create/Delete 内联完成，Read 本身保持简单的一次查询

**替代方案**:
- `time.Sleep(30 * time.Second)`：脆弱，违反项目 `helper.Retry` 惯例
- 完全不重试：会出现幽灵 drift

### Decision 5: Read 的一致性实现

**选择**: Read 调用 `service.DescribeVodAigcApiTokenById(ctx, subAppId, apiToken)`，内部执行 `DescribeAigcApiTokens` 并在返回的 `ApiTokens` 列表中线性查找当前 token 字符串；找到 → `_ = d.Set("api_token", apiToken)`；未找到 → `d.SetId("")`。

**理由**:
- 云侧 List API 一次性返回一个 SubAppId 下全部 token 字符串数组（API 文档未提及分页字段，SDK `DescribeAigcApiTokensRequestParams` 也仅有 `SubAppId`），故无分页负担
- 列表规模预期很小（每个应用通常维护个位数 Token）

**替代方案**:
- 缓存 Read 结果到 state：无必要且增加复杂度
- 调用单个 Token 查询接口：API 不提供

### Decision 6: service 层方法签名

在 `tencentcloud/services/vod/service_tencentcloud_vod.go` 中新增：

```go
// 返回云端生成的 apiToken
func (me *VodService) CreateVodAigcApiToken(ctx context.Context, subAppId uint64) (string, error)

// 查询指定 SubAppId 下全部 Token
func (me *VodService) DescribeVodAigcApiTokens(ctx context.Context, subAppId uint64) ([]string, error)

// 返回 token 是否仍存在（内部调用 DescribeVodAigcApiTokens 做线性查找）
func (me *VodService) DescribeVodAigcApiTokenById(ctx context.Context, subAppId uint64, apiToken string) (bool, error)

// 删除
func (me *VodService) DeleteVodAigcApiToken(ctx context.Context, subAppId uint64, apiToken string) error
```

所有方法内部使用 `me.client.UseVodClient().<Action>(request)` 调用 SDK，并走项目统一的 `LogElapsed` / `helper.Retry` 模式。

### Decision 7: 错误处理

**选择**: 统一使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装写操作，`tccommon.RetryError(e)` / `resource.NonRetryableError(e)` 区分瞬时/永久错误。

**特殊错误码处理**（基于 API 文档列出的错误码）:
- `ResourceNotFound.UserNotExist`：Delete 时应被视为"已删除"（幂等），不报错
- `FailedOperation.DBError`：可重试

### Decision 8: 敏感字段处理

**选择**: 
- schema 中 `api_token` 设 `Sensitive: true`
- service 层所有日志打印 request/response 前，不直接 `request.ToJsonString()` 暴露 token 值；改为仅打印 RequestId 与 SubAppId（Create 响应中的 Token 日志掩码为 `***`）
- Delete 日志同样不打印 token 明文

**理由**: 遵循通用秘钥凭证处理原则；`Sensitive: true` 可让 Terraform 自动在 CLI 输出打星号

### Decision 9: 文件命名与注册

- 实现: `tencentcloud/services/vod/resource_tc_vod_aigc_api_token.go`
- 测试: `tencentcloud/services/vod/resource_tc_vod_aigc_api_token_test.go`
- 文档模板: `tencentcloud/services/vod/resource_tc_vod_aigc_api_token.md`
- Provider 资源名: `tencentcloud_vod_aigc_api_token`
- 注册位置: `tencentcloud/provider.go` 的 `ResourcesMap`，以及 `tencentcloud/provider.md` VOD 分组

## Risks / Trade-offs

- **[风险] 30 秒同步延迟导致误判资源消失** → **缓解**: Create/Delete 后通过 `helper.Retry` 轮询 DescribeAigcApiTokens，最多等待 `ReadRetryTimeout`（3 分钟）
- **[风险] Token 明文泄漏** → **缓解**: schema `Sensitive: true`；service 层日志掩码；文档中提醒用户使用 `terraform output -raw` 时自担风险
- **[风险] Read 列表返回大数组** → **缓解**: 单个 SubAppId 下 Token 数量预期较少（业务上通常不会滥用），接受 O(n) 查找；若未来出现极端场景再引入本地缓存
- **[风险] Import 导入未知 token 失败** → **缓解**: Import 解析 `sub_app_id#api_token` 后走一次 Read，若列表中找不到则返回错误并提示用户核对 token
- **[风险] taint 重新创建会产生新 token，旧 token 可能残留云侧** → **缓解**: 正常 apply 流程是先 Destroy 再 Create（由 Terraform 框架自动处理）；若用户直接 `terraform state rm` 则确实会残留，属预期行为，文档中说明
- **[权衡] `api_token` 允许 `Optional` 以支持 import，但普通用户设置该字段可能引起混淆** → **接受**：文档示例仅保留 `sub_app_id`，并在文档中明确 import 用法

## Migration Plan

无迁移工作，这是一个全新的资源，不影响任何现有 TF 配置或 state。

部署步骤：
1. 合入代码后用户更新 Provider 版本即可使用
2. 用户在 `.tf` 文件中添加 `tencentcloud_vod_aigc_api_token` 资源块，指定 `sub_app_id`
3. `terraform apply` 创建，`terraform output <name>.api_token` 获取 token 值（标注 sensitive）

回滚: 移除资源块后 `terraform apply`，Provider 会自动调用 `DeleteAigcApiToken` 删除云侧 token。

## Open Questions

- `DescribeAigcApiTokens` 是否支持分页？  
  **当前决策**: 文档与 SDK `DescribeAigcApiTokensRequestParams` 均未声明 `Limit/Offset`，视为无分页；若后续 SDK 升级新增分页字段，再按项目现有模式内部自动分页
- 是否需要同步交付 `tencentcloud_vod_aigc_api_tokens` 数据源？  
  **当前决策**: 本 change 仅交付 resource；若用户后续提出数据源需求，单独立项
- `ResourceNotFound.UserNotExist` 触发条件是否包括"SubAppId 合法但无 token"？  
  **当前决策**: 不预设；实测过程中若发现则在 Read 分支中视为 `d.SetId("")`，而不是返回错误
