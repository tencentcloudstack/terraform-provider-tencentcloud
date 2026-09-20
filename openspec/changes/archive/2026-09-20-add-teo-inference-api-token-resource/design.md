## Context

腾讯云 EdgeOne（TEO）推理服务通过 Inference API Token 进行访问鉴权。当前 vendored 的 `tencentcloud-sdk-go`（`teov20220901` 包）已提供三个接口：

- `CreateInferenceAPIToken`：入参 `ZoneId`、`Name`，出参 `TokenId`、`Content`
- `DescribeInferenceAPITokens`：入参 `ZoneId`、`Offset`、`Limit`（最大 100），出参 `TotalCount`、`Tokens[]`（每个元素含 `TokenId`、`Name`、`Content`、`CreateTime`）
- `DeleteInferenceAPIToken`：入参 `ZoneId`、`TokenId`

这三个接口均为同步接口（SDK 中无异步标记），且**没有 Update 接口**。因此该 Terraform 资源属于 CRD 类型（仅 Create/Read/Delete），需要遵循 provider 中 CRD 资源的既有模式（如 `tencentcloud_teo_just_in_time_transcode_template`）。

约束：
- 复合 ID：使用 `zone_id` 与 `token_id` 以 `tccommon.FILED_SP` 分隔组成，支持 `terraform import`
- 无 Update 接口：`Id()` 字段设为 `ForceNew`，其余顶层可变字段加入 `immutableArgs`，变更时返回 error
- 资源代码风格严格参考 `tencentcloud_igtm_strategy` 资源（业务逻辑代码样式）
- 服务层 Describe 方法参考 `tencentcloud_teo_just_in_time_transcode_template` 中按 `token_id` 过滤查询的模式
- 单元测试使用 gomonkey mock 云 API（非 terraform 测试套件）

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_inference_api_token` 资源，完整支持 Token 的创建、查询、删除生命周期
- 资源使用 `zone_id#token_id` 联合 ID，支持 `terraform import`
- 通过 `immutableArgs` 机制阻止对不可变字段的更新，返回明确的 error 信息
- 通过 gomonkey 单元测试覆盖 Create/Read/Update(immutable 校验)/Delete 的业务逻辑
- 在 provider.go 中注册资源，并同步资源文档

**Non-Goals:**
- 不支持 Token 内容（`content`）的轮换或更新（云 API 无此能力）
- 不暴露 `Offset`/`Limit` 分页参数给用户（Read 内部使用最大 Limit 100 查询后按 `token_id` 过滤）
- 不新增 `_extension.go` 文件
- 不在资源 go 文件开头添加注释

## Decisions

### Decision 1: 资源类型为 CRD，使用 immutableArgs 模式

**选择**：参照 `tencentcloud_teo_just_in_time_transcode_template`，schema 中 `zone_id`、`name` 设为 `Required + ForceNew`；`token_id`、`content`、`create_time` 设为 `Computed`；Update 函数中声明 `immutableArgs := []string{"name"}`，检测到变更时返回 `fmt.Errorf`。

**备选**：不实现 Update 函数，让 Terraform 在检测到不可变字段变更时自动重建。

**理由**：
- 云 API 无 Update 接口，但有 Delete 接口，理论上可通过"删除+重建"实现更新。但 provider 现有 CRD 资源统一使用 `immutableArgs` 显式拦截，保持一致性
- `zone_id` 已是 `ForceNew`，重建时会触发 Delete+Create；`name` 加入 `immutableArgs` 后，若用户修改 `name`，Update 会先返回 error，由 Terraform 决定是否重建（因为 `name` 也是 `ForceNew`，实际上 Terraform 会直接重建，但保留 `immutableArgs` 校验符合既有规范且能覆盖非 ForceNew 字段场景）

### Decision 2: Read 通过 DescribeInferenceAPITokens 按 token_id 过滤

**选择**：`DescribeInferenceAPITokens` 接口不支持按 `TokenId` 直接过滤（入参只有 `ZoneId`、`Offset`、`Limit`），因此在 Read 中设置 `Limit = 100`（云 API 注释标注的最大值）从云端拉取 Token 列表，然后在客户端遍历 `Tokens` 数组匹配 `token_id`。

**备选**：分页循环拉取所有 Token。

**理由**：
- 该接口最大 Limit 为 100，单个站点的推理 Token 数量通常远小于此
- 参照 `DescribeJustInTimeTranscodeTemplates` 的实现，先查询列表再在客户端匹配，是该 provider 中缺少单资源查询接口时的标准做法
- 若未来 Token 数量超过 100，可再补充分页逻辑；当前以最大值 100 查询足够覆盖绝大多数场景，符合规则"若查询接口中有分页字段，则给定值应该是云API接口注释中标注的最大值"

### Decision 3: 复合 ID 格式 `zone_id#token_id`

**选择**：`d.SetId(zoneId + tccommon.FILED_SP + tokenId)`，Import 时用户需提供联合 ID。

**理由**：
- `zone_id` 与 `token_id` 共同唯一标识一个 Token，Delete 接口需要同时传入两者
- 与 `tencentcloud_teo_just_in_time_transcode_template`（`zone_id#template_id`）格式一致

### Decision 4: Read 接口返回空时的处理

**选择**：在 Read 的 retry 块内，若 `DescribeInferenceAPITokens` 返回的 `Tokens` 为空或未匹配到目标 `token_id`，返回 `NonRetryableError`；在 retry 失败路径上打印 `log.Printf("[CRUD] teo inference_api_token id=%s", d.Id())` 后 `d.SetId("")`。

**理由**：
- 遵循规则 #8：先打印日志保留现场，再 `d.SetId("")`，避免日志中无法定位是哪一次调用导致 id 被清空
- Token 不存在属于终态，不应让外层 retry 反复重试

### Decision 5: 单元测试使用 gomonkey mock 云 API

**选择**：参照 `resource_tc_teo_just_in_time_transcode_template_test.go` 与 `resource_tc_teo_create_cls_index_operation_test.go` 的模式，使用 gomonkey mock `UseTeoV20220901Client` 及对应的 Create/Describe/Delete 方法，复用 `newMockMeta`、`ptrString` 等已有的测试辅助函数。

**理由**：
- 规则要求新增资源使用 gomonkey mock 云 API 进行业务逻辑单元测试，不使用 terraform 测试套件
- teo 测试包中已存在共享的 `newMockMeta` 与 `ptrString` 等辅助函数，直接复用避免重复定义

## Risks / Trade-offs

- **Risk**：`DescribeInferenceAPITokens` 最大返回 100 条，若站点 Token 超过 100 则可能匹配不到目标 → **Mitigation**：当前以最大 Limit 100 查询；推理 Token 通常数量极少，实际场景不会触及上限；若触及可后续补充分页
- **Risk**：`content` 字段在 Create 时由云端返回并在 Read 时回填，属于敏感凭据信息 → **Mitigation**：schema 中标记为 `Computed` + `Sensitive: true`（与 SDK 中 Token 内容语义一致），避免明文输出
- **Trade-off**：CRD 资源修改 `name` 需要重建（删除旧 Token + 创建新 Token），会产生短暂的访问中断 → 可接受，与云 API 能力一致（无 Update 接口）