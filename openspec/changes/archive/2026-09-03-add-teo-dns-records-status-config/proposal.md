## Why

TEO（EdgeOne）DNS 记录的启用/停用状态目前无法通过 Terraform 声明式管理，用户只能通过控制台或 SDK 手动操作，导致运维流程断裂。腾讯云 TEO SDK 已提供 `DescribeDnsRecords`（查询 DNS 记录列表）和 `ModifyDnsRecordsStatus`（批量修改 DNS 记录状态）两个接口，可以新建一个 RESOURCE_KIND_CONFIG 资源 `tencentcloud_teo_dns_records_status`，在 Read 阶段通过 `DescribeDnsRecords` 查询记录状态，在 Update 阶段通过 `ModifyDnsRecordsStatus` 切换记录状态，实现声明式管理。

## What Changes

- 新增 Terraform RESOURCE_KIND_CONFIG 资源 `tencentcloud_teo_dns_records_status`，资源文件 `tencentcloud/services/teo/resource_tc_teo_dns_records_status_config.go`，只管理单个资源，不支持批量。
- 资源 Schema 仅包含 `ModifyDnsRecordsStatus` 接口的参数字段与状态控制字段：
  - `zone_id`（Required, ForceNew）：站点 ID，同时用于 Read（`DescribeDnsRecords` 的 `request.ZoneId`）和 Update（`ModifyDnsRecordsStatus` 的 `request.ZoneId`）。
  - `records_id`（Required, ForceNew）：DNS 记录 ID，与 `zone_id` 组合作为资源唯一 ID（`zone_id#records_id`）。
  - `status`（Required）：DNS 记录状态，枚举值 `enable`（已生效）/`disable`（已停用）。通过该字段控制记录的启用与停用，对应 `ModifyDnsRecordsStatus` 的 `RecordsToEnable`（status=enable 时）或 `RecordsToDisable`（status=disable 时）。
- 资源 ID 采用 `zone_id#records_id`（使用 `tccommon.FILED_SP` 即 `#` 作为分隔符），唯一标识一条 DNS 记录的状态配置。
- Create 操作：RESOURCE_KIND_CONFIG 资源无独立创建接口，Create 复用 Update 逻辑——仅 `d.SetId(zoneId + FILED_SP + recordsId)` 后直接调用 `resourceTencentCloudTeoDnsRecordsStatusUpdate(d, meta)`，与 `tencentcloud_teo_ddos_protection_config` 等其他 config 资源保持一致。
- Read 操作：调用 `DescribeDnsRecords`，传入 `zone_id` 与 `Filters`（按 `id` 过滤 `records_id`），查询对应 DNS 记录；若返回空则先打印 `log.Printf("[CRUD] teo_dns_records_status id=%s", d.Id())` 保留现场再 `d.SetId("")`，使用 `tccommon.ReadRetryTimeout` 重试。Read 回填 `status` 字段（取 `DnsRecord.Status`）。
- Update 操作：当 `status` 发生变化（`d.HasChange("status")`）时，调用 `ModifyDnsRecordsStatus`，根据 `status` 值将 `records_id` 放入 `RecordsToEnable`（enable）或 `RecordsToDisable`（disable），使用 `tccommon.WriteRetryTimeout` 重试。由于 `ModifyDnsRecordsStatus` 是异步接口，调用成功后轮询 `DescribeDnsRecords` 直到记录 `status` 达到目标值，再调用 `resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)` 回写最新状态。
- Delete 操作：CONFIG 资源无云端删除语义，Delete 为 no-op（不清空记录状态，仅移除 Terraform state）。
- 资源支持 import，import 时使用 `zone_id#records_id` 作为复合 ID。
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册资源 `tencentcloud_teo_dns_records_status`。
- 新增单元测试文件 `tencentcloud/services/teo/resource_tc_teo_dns_records_status_config_test.go`，使用 gomonkey mock 云 API，覆盖 Create/Read/Update 逻辑。
- 新增文档文件 `tencentcloud/services/teo/resource_tc_teo_dns_records_status.md`。

## Capabilities

### New Capabilities
- `teo-dns-records-status-config`: 新增 RESOURCE_KIND_CONFIG 资源 `tencentcloud_teo_dns_records_status`，通过 `DescribeDnsRecords` 读取 DNS 记录状态、通过 `ModifyDnsRecordsStatus` 更新记录启用/停用状态，只管理单个资源，使用 `zone_id#records_id` 作为唯一 ID，通过 `status` 字段控制记录状态，更新后轮询确认状态生效。

### Modified Capabilities
<!-- 无修改的已有 capability -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_dns_records_status_config.go`（新增资源实现）
  - `tencentcloud/services/teo/resource_tc_teo_dns_records_status_config_test.go`（新增单元测试，gomonkey mock）
  - `tencentcloud/services/teo/resource_tc_teo_dns_records_status.md`（新增资源文档）
  - `tencentcloud/provider.go`（注册资源）
  - `tencentcloud/provider.md`（注册资源）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.DescribeDnsRecordsRequest`、`teov20220901.ModifyDnsRecordsStatusRequest`、`teov20220901.DnsRecord`、`teov20220901.AdvancedFilter`，无需变更 vendor。
- 向后兼容：新增资源，不影响已有资源与 state。
- 文档：需要通过 `make doc` 生成 `website/docs/` 下的 markdown 文档（由收尾阶段 tfpacer-finalize skill 执行）。
