## Context

TEO（EdgeOne）DNS 记录的启用/停用状态目前无法通过 Terraform 声明式管理。腾讯云 TEO SDK（`teov20220901`）已提供两个相关接口：
- `DescribeDnsRecords`：查询 DNS 记录列表，入参包含 `ZoneId`、`Offset`、`Limit`、`Filters`（`AdvancedFilter`，含 `Name`/`Values`/`Fuzzy`）、`SortBy`、`SortOrder`、`Match`，出参为 `DnsRecords`（`DnsRecord` 列表，含 `ZoneId`/`RecordId`/`Name`/`Type`/`Location`/`Content`/`TTL`/`Weight`/`Priority`/`Status`/`CreatedOn`/`ModifiedOn`）。`Filters` 支持 `id` 过滤条件，可按 DNS 记录 ID 精确/模糊查询。
- `ModifyDnsRecordsStatus`：批量修改 DNS 记录状态（**异步接口**），入参包含 `ZoneId`、`RecordsToEnable`（待启用记录 ID 列表，上限 200）、`RecordsToDisable`（待停用记录 ID 列表，上限 200），同一个记录 ID 不能同时存在于两个列表中。

当前状态：
- 不存在 `tencentcloud_teo_dns_records_status` 资源，provider.go 中无注册。
- SDK 中 `DnsRecord`、`AdvancedFilter` 已 vendored 可用。

约束：
- 本次为 RESOURCE_KIND_CONFIG 资源：只管理单个资源，不支持批量；资源存在配置就存在，主要管理配置的读取和更新，需要创建 RU 接口。
- 资源文件命名格式：`resource_tc_teo_dns_records_status_config.go`。
- 参照其他 config 类型资源做重构，schema 字段只需要 `zone_id`、`records_id`、`status`，不需要暴露 Describe 接口的查询参数（`filters`、`sort_by`、`sort_order`、`match`）和返回的列表数据（`dns_records`）。
- 必须保持向后兼容（新增资源，无破坏性）。
- 调用云 API 需以 `tccommon.ReadRetryTimeout`（Read）/`tccommon.WriteRetryTimeout`（Update）作为超时时间添加 retry 处理，失败时用 `tccommon.RetryError()` 包装。
- Read 回填前需判断 Response 字段是否为 nil，nil 则不调用 set；若云 API 返回空，先打印 `log.Printf("[CRUD] xxx id=%s", d.Id())` 保留现场再 `d.SetId("")`。
- Create 复用 Update 的简化模式：`d.SetId(zoneId + FILED_SP + recordsId)` 后直接 `return Update(d, meta)`，不在 Create 中重复调用 API 的逻辑。
- 只管理单个资源，不支持批量：`records_id` 为单个记录 ID，通过 `status` 控制其启用/停用。
- 资源支持 import，import 时使用 `zone_id#records_id` 作为复合 ID。
- `ModifyDnsRecordsStatus` 是异步接口，更新后需轮询 `DescribeDnsRecords` 确认记录 `status` 达到目标值。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_dns_records_status` 资源，通过 `DescribeDnsRecords` 读取 DNS 记录状态，通过 `ModifyDnsRecordsStatus` 更新记录启用/停用状态。
- 只管理单个资源，使用 `records_id`（单个记录 ID）+ `status`（enable/disable）控制记录状态。
- Schema 只包含 `zone_id`、`records_id`、`status` 三个字段。
- 资源 ID 使用 `zone_id#records_id` 复合 ID。
- Update 阶段调用 `ModifyDnsRecordsStatus` 后轮询 `DescribeDnsRecords` 确认状态生效（异步接口）。
- 资源支持 import（使用复合 ID）。
- 单元测试使用 gomonkey mock 云 API，覆盖 Create/Read/Update。

**Non-Goals:**
- 不修改已有 `tencentcloud_teo_dns_record` 资源的任何行为。
- 不支持批量管理多条 DNS 记录状态（用户自定义要求只管理单个资源）。
- 不暴露 `filters`、`sort_by`、`sort_order`、`match` 等 Describe 接口查询参数。
- 不暴露 `dns_records` 列表数据到 schema。
- 不新增 `_extension.go` 文件。
- 不手动编写 `website/docs/` 文档（由 `make doc` 自动生成，收尾阶段执行）。

## Decisions

### Decision 1: 资源 ID 使用 `zone_id#records_id` 复合 ID

**选择**：资源 ID 使用 `zone_id` + `tccommon.FILED_SP`（即 `#`）+ `records_id` 作为复合 ID，唯一标识一条 DNS 记录的状态配置。

**备选**：仅用 `zone_id` 作为 ID。

**理由**：
- 用户明确要求"zone_id#records_id 是唯一 id"。
- 一个 `zone_id` 下存在多条 DNS 记录，仅用 `zone_id` 无法区分具体管理哪条记录的状态。
- 使用复合 ID 后，可对特定记录进行精确的状态读取与更新。
- import 时使用 `zone_id#records_id` 作为复合 ID，明确定位资源。

### Decision 2: Create 复用 Update 逻辑（简化模式）

**选择**：RESOURCE_KIND_CONFIG 无独立创建接口，Create 函数采用与其他 config 资源一致的简化模式：仅 `d.SetId(zoneId + FILED_SP + recordsId)` 后直接 `return resourceTencentCloudTeoDnsRecordsStatusUpdate(d, meta)`，不在 Create 中重复调用 `ModifyDnsRecordsStatus` 的逻辑。

**备选**：Create 中独立构造 `ModifyDnsRecordsStatusRequest` 并调用 API，再单独轮询状态。

**理由**：
- CONFIG 资源语义为"资源存在配置就存在"，Create 等同于首次设置配置，Update 已经包含了构造请求、调用 API、retry、轮询状态、调用 Read 回填的全部逻辑。
- 复用 Update 避免代码重复，与仓库内既有 config 资源实现风格一致，便于维护。
- Create 设置 ID 后，Update 内通过 `d.HasChange("status")` 判断字段是否变化来决定是否调用 API；对于全新资源，旧状态为空，设置了 `status` 即视为变化，会触发调用。

### Decision 3: Schema 使用 `records_id` + `status` 替代 `records_to_enable`/`records_to_disable`

**选择**：schema 字段为 `zone_id`（Required, ForceNew）、`records_id`（Required, ForceNew）、`status`（Required，枚举 `enable`/`disable`）。不使用 `records_to_enable`/`records_to_disable` 列表字段，改用单个 `records_id` + `status` 控制记录状态。

**备选**：保留 `records_to_enable`/`records_to_disable` 列表字段。

**理由**：
- 用户明确要求"不需要 records_to_enable 和 records_to_disable，只需要一个 records_id 即可，再增加 status 字段"。
- `status` 字段语义清晰，枚举 `enable`（已生效）/`disable`（已停用），与 `DnsRecord.Status` 出参值一致。
- 单个 `records_id` + `status` 比"两个列表字段"更符合声明式模型，避免用户在两个互斥列表间维护一致性。
- Update 时根据 `status` 值将 `records_id` 放入 `RecordsToEnable` 或 `RecordsToDisable`，映射到云 API 入参。

### Decision 4: Read 使用 Filters 按 records_id 过滤查询

**选择**：Read 调用 `DescribeDnsRecords` 传入 `zone_id` 与 `Filters`（`Name="id"`, `Values=[records_id]`），查询对应 DNS 记录；若返回空则先打印日志保留现场再 `d.SetId("")`。Read 回填 `status` 字段（取 `DnsRecord.Status`）。

**理由**：
- 复合 ID 中包含 `records_id`，Read 需精确定位该记录。
- `DescribeDnsRecords` 的 `Filters` 支持 `id` 过滤条件，可按记录 ID 查询。
- Read 回填 `status` 使 Terraform state 与云端实际状态保持一致。

### Decision 5: Update 后轮询确认状态生效（异步接口）

**选择**：`ModifyDnsRecordsStatus` 是异步接口，调用成功（retry 块外）后，使用 `resource.StateChangeConf` 轮询 `DescribeDnsRecords` 直到记录 `status` 达到目标值（`enable` 或 `disable`），再调用 `resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)` 回写最新状态。

**备选**：调用成功后直接调用 Read 回填，不轮询。

**理由**：
- 用户明确要求"更新接口是异步接口，查询接口要轮询 records 状态"。
- 异步接口调用返回成功不代表状态已生效，需轮询确认，避免 Read 拿到中间态导致 plan drift。
- 使用 `resource.StateChangeConf`（`Delay`、`MinTimeout`、`Timeout=tccommon.ReadRetryTimeout`）实现轮询，与仓库内 `tencentcloud_teo_alias_domain` 等资源处理异步操作的方式一致。

### Decision 6: Delete 为 no-op

**选择**：Delete 函数为 no-op，不清空记录状态，仅移除 Terraform state。

**理由**：
- CONFIG 资源无云端删除语义，DNS 记录本身由 `tencentcloud_teo_dns_record` 资源管理生命周期，本资源只管理状态配置。
- 销毁时不回滚状态，避免误操作影响线上解析。

## Risks / Trade-offs

- **Risk**：`ModifyDnsRecordsStatus` 异步调用后轮询超时未达到目标状态 → **Mitigation**：使用 `tccommon.ReadRetryTimeout` 作为轮询超时，超时后返回错误提示人工介入。
- **Trade-off**：`status` 为 Required 且非 ForceNew，变更 `status` 触发 Update 调用 `ModifyDnsRecordsStatus`，符合 CONFIG 资源"主要管理配置更新"的语义。
- **Risk**：CONFIG 资源 Create 复用 Update，若 `status` 未变化则 `d.HasChange` 为 false，Update 不调用 API，仅调用 Read → **Mitigation**：对于全新资源，旧状态为空，设置了 `status` 即视为变化，会触发调用。

## Migration Plan

- 新增资源为纯加法，无 state 迁移需求。
- 存量资源：无（全新资源）。
- 文档更新：新增 `resource_tc_teo_dns_records_status.md`，通过 `make doc` 生成 `website/docs/` 文档（收尾阶段执行）。
- 回滚：若需要回退，移除 provider.go 注册与资源文件即可，无 state 残留。

## Open Questions

- 无
