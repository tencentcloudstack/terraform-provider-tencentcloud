## Why

TEO（EdgeOne）DNS 记录的启用/停用状态目前无法通过 Terraform 声明式管理，用户只能通过控制台或 SDK 手动操作，导致运维流程断裂。腾讯云 TEO SDK 已提供 `DescribeDnsRecords`（查询 DNS 记录列表）和 `ModifyDnsRecordsStatus`（批量修改 DNS 记录状态）两个接口，可以新建一个 RESOURCE_KIND_CONFIG 资源 `tencentcloud_teo_dns_records_status`，在 Read 阶段通过 `DescribeDnsRecords` 查询记录并回填状态，在 Update 阶段通过 `ModifyDnsRecordsStatus` 切换记录状态，实现声明式管理。

## What Changes

- 新增 Terraform RESOURCE_KIND_CONFIG 资源 `tencentcloud_teo_dns_records_status`，资源文件 `tencentcloud/services/teo/resource_tc_teo_dns_records_status_config.go`，只管理单个资源，不支持批量。
- 资源 Schema 包含：
  - `zone_id`（Required, ForceNew）：站点 ID，同时用于 Read（`DescribeDnsRecords` 的 `request.ZoneId`）和 Update（`ModifyDnsRecordsStatus` 的 `request.ZoneId`）。
  - `filters`（Optional）：过滤条件列表，元素为 `AdvancedFilter`，子字段 `name`（Required）、`values`（Required）、`fuzzy`（Optional），用于 Read 时按条件查询 DNS 记录。
  - `sort_by`（Optional）：排序依据，对应 `request.SortBy`。
  - `sort_order`（Optional）：排序方式，对应 `request.SortOrder`。
  - `match`（Optional）：匹配方式，对应 `request.Match`。
  - `records_to_enable`（Optional）：待启用的 DNS 记录 ID 列表，对应 `request.RecordsToEnable`，Update 时传入。
  - `records_to_disable`（Optional）：待停用的 DNS 记录 ID 列表，对应 `request.RecordsToDisable`，Update 时传入。
  - `dns_records`（Computed）：DNS 记录列表，对应 `response.Response.DnsRecords`，展开为列表，每个元素包含 `zone_id`、`record_id`、`name`、`type`、`location`、`content`、`ttl`、`weight`、`priority`、`status`、`created_on`、`modified_on` 字段。
- 资源 ID 采用复合 ID：`zone_id`（CONFIG 资源只管理单个资源的配置，ID 仅用 `zone_id` 标识）。
- Create 操作：RESOURCE_KIND_CONFIG 资源无独立创建接口，Create 复用 Update 逻辑（调用 `ModifyDnsRecordsStatus` 设置初始状态），成功后调用 Read 回填。
- Read 操作：调用 `DescribeDnsRecords`，传入 `zone_id`、`filters`（若有）、`sort_by`、`sort_order`、`match`，查询 DNS 记录列表并回填 `dns_records`（取列表第一项的各字段平铺到顶层 schema），使用 `tccommon.ReadRetryTimeout` 重试。
- Update 操作：调用 `ModifyDnsRecordsStatus`，传入 `zone_id`、`records_to_enable`（若有）、`records_to_disable`（若有），使用 `tccommon.WriteRetryTimeout` 重试；成功后调用 Read 轮询直到接口生效（检查 `dns_records` 中对应记录的 `status` 达到期望值）。
- Delete 操作：CONFIG 资源无云端删除语义，Delete 为 no-op（不清空记录状态，仅移除 Terraform state）。
- 资源支持 import，import 时使用 `zone_id` 作为 ID。
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册资源 `tencentcloud_teo_dns_records_status`。
- 新增单元测试文件 `tencentcloud/services/teo/resource_tc_teo_dns_records_status_config_test.go`，使用 gomonkey mock 云 API，覆盖 Create/Read/Update 逻辑。
- 新增文档文件 `tencentcloud/services/teo/resource_tc_teo_dns_records_status.md`。

## Capabilities

### New Capabilities
- `teo-dns-records-status-config`: 新增 RESOURCE_KIND_CONFIG 资源 `tencentcloud_teo_dns_records_status`，通过 `DescribeDnsRecords` 读取 DNS 记录状态、通过 `ModifyDnsRecordsStatus` 更新记录启用/停用状态，只管理单个资源。

### Modified Capabilities
<!-- 无修改的已有 capability -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_dns_records_status_config.go`（新增资源实现）
  - `tencentcloud/services/teo/resource_tc_teo_dns_records_status_config_test.go`（新增单元测试，gomonkey mock）
  - `tencentcloud/services/teo/resource_tc_teo_dns_records_status.md`（新增资源文档）
  - `tencentcloud/provider.go`（注册资源）
  - `tencentcloud/provider.md`（注册资源）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.DescribeDnsRecordsRequest`、`teov20220901.ModifyDnsRecordsStatusRequest`、`teov20220901.AdvancedFilter`、`teov20220901.DnsRecord`，无需变更 vendor。
- 向后兼容：新增资源，不影响已有资源与 state。
- 文档：需要通过 `make doc` 生成 `website/docs/` 下的 markdown 文档（由收尾阶段 tfpacer-finalize skill 执行）。
