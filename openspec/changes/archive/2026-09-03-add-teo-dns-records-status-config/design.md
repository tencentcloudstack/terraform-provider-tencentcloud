## Context

TEO（EdgeOne）DNS 记录的启用/停用状态目前无法通过 Terraform 声明式管理。腾讯云 TEO SDK（`teov20220901`）已提供两个相关接口：
- `DescribeDnsRecords`：查询 DNS 记录列表，入参包含 `ZoneId`、`Offset`、`Limit`、`Filters`（`AdvancedFilter`，含 `Name`/`Values`/`Fuzzy`）、`SortBy`、`SortOrder`、`Match`，出参为 `DnsRecords`（`DnsRecord` 列表，含 `ZoneId`/`RecordId`/`Name`/`Type`/`Location`/`Content`/`TTL`/`Weight`/`Priority`/`Status`/`CreatedOn`/`ModifiedOn`）。
- `ModifyDnsRecordsStatus`：批量修改 DNS 记录状态，入参包含 `ZoneId`、`RecordsToEnable`（待启用记录 ID 列表，上限 200）、`RecordsToDisable`（待停用记录 ID 列表，上限 200），同一个记录 ID 不能同时存在于两个列表中。

当前状态：
- 不存在 `tencentcloud_teo_dns_records_status` 资源，provider.go 中无注册。
- SDK 中 `DnsRecord` 已 vendored 可用。

约束：
- 本次为 RESOURCE_KIND_CONFIG 资源：只管理单个资源，不支持批量；资源存在配置就存在，主要管理配置的读取和更新，需要创建 RU 接口。
- 资源文件命名格式：`resource_tc_teo_dns_records_status_config.go`。
- 参照其他 config 类型资源做重构，schema 字段只需要 Modify 接口的参数字段（`zone_id`、`records_to_enable`、`records_to_disable`），不需要暴露 Describe 接口的查询参数（`filters`、`sort_by`、`sort_order`、`match`）和返回的列表数据（`dns_records`）。
- 必须保持向后兼容（新增资源，无破坏性）。
- 调用云 API 需以 `tccommon.ReadRetryTimeout`（Read）/`tccommon.WriteRetryTimeout`（Update）作为超时时间添加 retry 处理，失败时用 `tccommon.RetryError()` 包装。
- `ModifyDnsRecordsStatus` 为异步接口，调用后需调用 `DescribeDnsRecords` 轮询直到接口生效（记录 `status` 达到 enable/disable 期望值）。
- Read 回填前需判断 Response 字段是否为 nil，nil 则不调用 set；若云 API 返回空，先打印 `log.Printf("[CRUD] xxx id=%s", d.Id())` 保留现场再 `d.SetId("")`。
- Create 完成后必须检查返回值是否为空，若为空返回 `NonRetryableError`。
- 只管理单个资源，不支持批量：`records_to_enable`/`records_to_disable` 传入单个记录 ID。
- 资源支持 import，import 时使用 `zone_id` 作为 ID。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_dns_records_status` 资源，通过 `DescribeDnsRecords` 读取 DNS 记录状态，通过 `ModifyDnsRecordsStatus` 更新记录启用/停用状态。
- 只管理单个资源，`records_to_enable`/`records_to_disable` 传入单个记录 ID。
- Schema 只包含 Modify 接口参数字段：`zone_id`、`records_to_enable`、`records_to_disable`。
- Update 阶段调用 `ModifyDnsRecordsStatus` 后轮询 `DescribeDnsRecords` 确认状态生效。
- 资源支持 import。
- 单元测试使用 gomonkey mock 云 API，覆盖 Create/Read/Update。

**Non-Goals:**
- 不修改已有 `tencentcloud_teo_dns_record` 资源的任何行为。
- 不支持批量管理多条 DNS 记录状态（用户自定义要求只管理单个资源）。
- 不暴露 `filters`、`sort_by`、`sort_order`、`match` 等 Describe 接口查询参数。
- 不暴露 `dns_records` 列表数据到 schema。
- 不新增 `_extension.go` 文件。
- 不手动编写 `website/docs/` 文档（由 `make doc` 自动生成，收尾阶段执行）。

## Decisions

### Decision 1: 资源 ID 使用 `zone_id`

**选择**：资源 ID 仅用 `zone_id` 标识，不拼接 record_id。

**备选**：使用 `zone_id#record_id` 复合 ID。

**理由**：
- RESOURCE_KIND_CONFIG 资源描述一个资源的配置，资源存在配置就存在；本资源以 `zone_id` 为配置归属。
- 用户自定义要求"只管理单个资源，不支持批量"，配置入口为 `zone_id`。
- import 时使用 `zone_id` 作为 ID，简洁明了。

### Decision 2: Create 复用 Update 逻辑

**选择**：RESOURCE_KIND_CONFIG 无独立创建接口，Create 函数复用 Update 逻辑：调用 `ModifyDnsRecordsStatus` 设置初始状态，成功后调用 Read 回填。

**备选**：Create 仅调用 Read 不调用 Modify。

**理由**：
- CONFIG 资源语义为"资源存在配置就存在"，Create 等同于首次设置配置。
- 用户提供了 `ModifyDnsRecordsStatus` 作为 Update 接口，Create 复用该接口设置 `records_to_enable`/`records_to_disable` 初始值，符合 CONFIG 资源 RU 模式。

### Decision 3: Schema 只包含 Modify 接口参数字段

**选择**：参照其他 config 类型资源做重构，schema 只包含 `ModifyDnsRecordsStatus` 接口的参数字段：`zone_id`（Required, ForceNew）、`records_to_enable`（Optional）、`records_to_disable`（Optional）。不暴露 `filters`、`sort_by`、`sort_order`、`match` 等 Describe 接口查询参数，也不暴露 `dns_records` 列表数据。

**备选**：暴露 Describe 接口的查询参数和返回数据。

**理由**：
- 用户明确要求"参照其他的config类型资源做重构，schema字段只需要修改接口的就行"。
- CONFIG 资源核心是配置的更新与读取，schema 字段应与 Modify 接口参数一一对应。
- `filters`/`sort_by`/`sort_order`/`match` 仅用于 Describe 查询，不属于配置本身。
- `dns_records` 是查询结果而非配置项，不暴露到 schema 可简化资源模型。

### Decision 4: Read 仅用 zone_id 查询判断资源是否存在

**选择**：Read 调用 `DescribeDnsRecords` 仅传入 `zone_id`，查询 DNS 记录列表判断资源是否存在；若返回空则 `d.SetId("")`，不回填任何列表数据到 schema。

**理由**：
- Schema 不含 `dns_records` 字段，Read 只需判断资源是否存在。
- 若 `zone_id` 下存在 DNS 记录，则配置存在；否则配置不存在。

### Decision 5: Update 后轮询 DescribeDnsRecords 确认状态生效

**选择**：`ModifyDnsRecordsStatus` 为异步接口，调用成功后调用 `DescribeDnsRecords` 轮询，检查目标记录 `status` 是否达到期望值（enable/disable），在 `d.Timeout(schema.TimeoutUpdate)` 内等待。

**备选**：不等待，立即返回。

**理由**：
- 异步接口需轮询确认生效，避免 apply 结束后 Read 拿到中间态导致 plan drift。
- 与 provider 现有异步操作处理模式一致。

### Decision 6: Delete 为 no-op

**选择**：Delete 函数为 no-op，不清空记录状态，仅移除 Terraform state。

**理由**：
- CONFIG 资源无云端删除语义，DNS 记录本身由 `tencentcloud_teo_dns_record` 资源管理生命周期，本资源只管理状态配置。
- 销毁时不回滚状态，避免误操作影响线上解析。

### Decision 7: 只管理单个资源，records_to_enable/records_to_disable 传单个 ID

**选择**：`records_to_enable`/`records_to_disable` schema 为 `TypeList` of `TypeString`（与云 API 入参类型一致），但用户语义上只管理单个资源，使用时传入单个记录 ID。

**理由**：
- 云 API 入参 `RecordsToEnable`/`RecordsToDisable` 为 `[]*string`，schema 保持 `TypeList` 以匹配 API 结构。
- 用户自定义要求"只管理单个资源，不支持批量"，文档示例展示单个 ID 用法。

## Risks / Trade-offs

- **Risk**：`ModifyDnsRecordsStatus` 异步轮询超时 → **Mitigation**：沿用 provider 统一 `RetryError` 模式，用户可重跑 apply 收敛。
- **Trade-off**：`records_to_enable`/`records_to_disable` 为 `TypeList` 但语义上只传单个 ID，可能与"只管理单个资源"语义略有张力 → 可接受，schema 类型与云 API 一致，文档约束用法。
- **Risk**：CONFIG 资源 Create 复用 Update 调用 `ModifyDnsRecordsStatus`，若用户未配置 `records_to_enable`/`records_to_disable` 则不调用 API → **Mitigation**：Create 时检查两个字段是否非空，都为空则跳过 Modify 直接 Read。

## Migration Plan

- 新增资源为纯加法，无 state 迁移需求。
- 存量资源：无（全新资源）。
- 文档更新：新增 `resource_tc_teo_dns_records_status.md`，通过 `make doc` 生成 `website/docs/` 文档（收尾阶段执行）。
- 回滚：若需要回退，移除 provider.go 注册与资源文件即可，无 state 残留。

## Open Questions

- 无
