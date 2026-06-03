## Context

`tencentcloud_dnspod_record` 资源底层使用 DNSPod SDK `DescribeRecord` 接口（`v20210323`），返回 `RecordInfo` 结构体已包含 `UpdatedOn`（记录最后更新时间）字段（vendor `models.go:7523`）。`tencentcloud_dnspod_record_list` 数据源底层使用 `DescribeRecordList`，返回 `RecordListItem.UpdatedOn`（vendor `models.go:7541`）。

当前状态：
- 资源：schema 中**无** `updated_on`，Read 中**未读取** `recordInfo.UpdatedOn`。
- 数据源：schema 中 `record_list` 与 `instance_list` 子块**已包含** `updated_on`（string、Computed），Read 中**已正确赋值** `recordListItemMap["updated_on"] = recordListItem.UpdatedOn` / `instanceListItemMap["updated_on"] = recordListItem.UpdatedOn`。

本变更主要是给资源补齐 Read 回写，以及把数据源既有字段纳入正式 spec 承诺。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_dnspod_record` 资源中暴露 `updated_on`（Computed string），Read 阶段从 `RecordInfo.UpdatedOn` 回写。
- 在 spec 中正式声明数据源 `record_list[].updated_on` 与 `instance_list[].updated_on` 为 Computed string。
- 文档/样例同步更新；通过 `make doc` 重新生成 website 文档。

**Non-Goals:**
- 不修改任何现有字段语义，不动 Create/Update/Delete 逻辑。
- 不引入分页/排序/过滤的新参数。
- 不实现"基于 `updated_on` 触发 Update"等业务联动逻辑（云端字段为只读时间戳）。
- 不为 `created_on` 等其他时间戳一并新增（聚焦 `updated_on`，避免范围蔓延；后续可单独提案）。

## Decisions

### D1：字段名使用 `updated_on`，类型 `schema.TypeString`

- **选择**：snake_case 风格 `updated_on`，与项目惯例（如 `created_at` / `updated_at` / 已有数据源 `updated_on`）一致；类型 `string`，直接透传 SDK 返回的字符串时间格式（云端通常返回 `YYYY-MM-DD HH:MM:SS`）。
- **替代**：
  - `last_updated_on` —— 与 SDK 字段命名偏离，且数据源已用 `updated_on`，应保持一致。
  - `TypeInt`（unix 时间戳）—— SDK 返回 string，需要解析转换，可能产生时区/解析问题；保持原样最安全。

### D2：仅在资源端新增 schema + Read，数据源不改代码

- **选择**：
  - `resource_tc_dnspod_record.go`：schema 新增 Computed 字段；Read 中 `if recordInfo.UpdatedOn != nil { _ = d.Set("updated_on", recordInfo.UpdatedOn) }`。
  - `data_source_tc_dnspod_record_list.go`：仅同步文档与 spec，不改代码（schema 与 Read 已具备字段）。
- **理由**：数据源现状代码与目标行为一致；改代码反而引入风险。Spec 只需 ADDED Requirement 描述资源端新行为 + 数据源端既有承诺。

### D3：使用 nil 检查 + `d.Set` 单值（与现有字段同风格）

- 与同文件内 `recordInfo.SubDomain` / `recordInfo.MX` 等字段保持同一回写风格：`_ = d.Set("updated_on", recordInfo.UpdatedOn)`，对 nil 做防护。

### D4：测试断言使用 `TestCheckResourceAttrSet` 而非具体值

- `updated_on` 为云端时间戳，每次创建结果不同，无法硬编码。使用 `TestCheckResourceAttrSet("tencentcloud_dnspod_record.foo","updated_on")` 验证非空即可。

### D5：仅新增 Computed 字段，向后兼容

- 不影响存量 state（state 中 `updated_on` 不存在 → refresh 后填入），不需要 `state migration`。

## Risks / Trade-offs

- **[风险] SDK 在某些场景返回 `UpdatedOn=nil`** → 已加 nil 检查，Set 跳过，state 字段保持空字符串/未设置；不会破坏 plan。
- **[风险] 时间格式可能因地域差异返回不同形式** → 直接透传字符串，不在 Provider 侧解析；用户层若需要解析自行处理。
- **[Trade-off] 不在本次同时新增 `created_on`** → 聚焦最小可审，未来可独立加；Capability 命名 `dnspod-record-updated-on` 已明确边界。

## Migration Plan

无破坏性变更：
1. 升级 Provider 后 `terraform refresh` / 下次 `terraform plan` 自动把 `updated_on` 填入 state。
2. 用户已有 HCL 无需变更；如需引用，使用 `tencentcloud_dnspod_record.X.updated_on`。
3. 回滚：直接降级 Provider 即可，state 中多出的 `updated_on` 字段在旧版本中会被忽略（schema 不存在则 plan/apply 不读不写）。

## Open Questions

无。
