## ADDED Requirements

### Requirement: REQ-DNSPOD-RECORD-UPDATED-ON-001 - dnspod_record 资源暴露 updated_on Computed 字段

The `tencentcloud_dnspod_record` resource MUST expose a Computed `updated_on` attribute (string) that mirrors the `UpdatedOn` field returned by the DNSPod `DescribeRecord` API (`RecordInfo.UpdatedOn`).

`tencentcloud_dnspod_record` 资源必须新增 Computed 类型为 string 的 `updated_on` 属性，并在 Read 阶段将 DNSPod `DescribeRecord` 接口返回的 `RecordInfo.UpdatedOn` 透传写入 Terraform state，使用户能在 `terraform output` / 跨资源引用中访问该记录的最后更新时间。

**约束**：

- 类型为 `schema.TypeString`、`Computed: true`，**不**设置 `Optional`/`Required`/`ForceNew`/`Default`。
- 仅 Read 阶段写入；Create/Update/Delete 不感知该字段。
- 当 SDK 返回 `RecordInfo.UpdatedOn == nil` 时，Provider MUST 跳过 `d.Set`，state 中保持原值不被错误清空。
- 仅新增字段，不修改任何已有字段语义；向后兼容。

#### Scenario: 资源 Read 回写 updated_on

- **WHEN** 用户对一条已存在的 `tencentcloud_dnspod_record` 执行 `terraform refresh` 或 `terraform plan`
- **THEN**
  - Provider 调用 `DescribeRecord` 接口，得到 `RecordInfo.UpdatedOn`（如 `"2026-06-03 09:30:00"`）
  - 当 `RecordInfo.UpdatedOn != nil` 时，Provider 通过 `d.Set("updated_on", recordInfo.UpdatedOn)` 写入 state
  - `terraform state show <addr>` 输出 `updated_on = "2026-06-03 09:30:00"`
  - 后续 `terraform plan` 不产生意料之外的 diff

**验收标准**：
- ✅ schema 中存在 `updated_on`：`Type: schema.TypeString, Computed: true`
- ✅ Read 函数中存在 `_ = d.Set("updated_on", recordInfo.UpdatedOn)` 调用
- ✅ `TestCheckResourceAttrSet("tencentcloud_dnspod_record.<name>", "updated_on")` 在验收测试中通过

---

#### Scenario: 创建后立即可读 updated_on

- **WHEN** 用户执行 `terraform apply` 创建一条新的 `tencentcloud_dnspod_record`
- **THEN**
  - Create 完成后调用 Read，state 中 `updated_on` 字段被首次填充
  - `terraform output` 或在其它资源中通过 `tencentcloud_dnspod_record.X.updated_on` 引用可获取非空字符串

**验收标准**：
- ✅ 创建后无需手动 refresh，`terraform state show` 已含 `updated_on`
- ✅ `output` 块输出该值非空

---

#### Scenario: SDK 返回 UpdatedOn 为 nil 的容错

- **WHEN** SDK 因边缘场景返回 `RecordInfo.UpdatedOn == nil`
- **THEN**
  - Provider MUST 跳过 `d.Set("updated_on", ...)` 调用
  - state 中已有 `updated_on` 值不被错误清空为空字符串
  - 不抛出错误，Read 流程正常完成

**验收标准**：
- ✅ Read 函数在赋值前进行 `if recordInfo.UpdatedOn != nil` 防护
- ✅ 单元/手工触发 nil 路径时 Read 返回 `nil` error

---

#### Scenario: 旧 state 升级兼容

- **WHEN** 用户从未含 `updated_on` 字段的旧 Provider 版本升级到本版本，对存量资源执行 `terraform plan`
- **THEN**
  - 不要求 state migration
  - `updated_on` 作为 Computed 字段，在下一次 Read 时被自动填充
  - 不产生破坏性变更或额外 diff（除该 Computed 字段从空到有值）

**验收标准**：
- ✅ 升级后存量 HCL 无需修改
- ✅ 不需要 state migration 步骤

---

### Requirement: REQ-DNSPOD-RECORD-UPDATED-ON-002 - dnspod_record_list 数据源暴露嵌套 updated_on Computed 字段

The `tencentcloud_dnspod_record_list` data source MUST expose Computed `updated_on` attributes (string) inside both the `record_list` and `instance_list` nested blocks, mapped to `RecordListItem.UpdatedOn` returned by the DNSPod `DescribeRecordList` API.

`tencentcloud_dnspod_record_list` 数据源必须在 `record_list` 与 `instance_list` 嵌套块下分别暴露 Computed string 类型的 `updated_on` 字段，并在 Read 阶段将 SDK `RecordListItem.UpdatedOn` 透传写入。

**约束**：

- `record_list[].updated_on`：`schema.TypeString, Computed: true`
- `instance_list[].updated_on`：`schema.TypeString, Computed: true`
- 当某条 `RecordListItem.UpdatedOn == nil` 时，Read 跳过该条目的 `updated_on` 赋值，不污染其它字段。
- 与现有 schema/Read 行为一致，本需求是对既有实现的形式化承诺，不引入代码变更。

#### Scenario: 数据源返回包含 updated_on 的 record_list

- **WHEN** 用户在 HCL 中使用 `data "tencentcloud_dnspod_record_list" "foo" { domain = "..." }` 并执行 `terraform plan/apply`
- **THEN**
  - `data.tencentcloud_dnspod_record_list.foo.record_list[*].updated_on` 可被引用
  - 对至少 1 条云端返回有 `UpdatedOn` 的记录，对应数组项的 `updated_on` 为非空字符串
  - `instance_list[*].updated_on` 同步具备一致语义

**验收标准**：
- ✅ schema 中两处嵌套块都包含 `updated_on`：`Type: schema.TypeString, Computed: true`
- ✅ Read 中存在 `recordListItemMap["updated_on"] = recordListItem.UpdatedOn` 与 `instanceListItemMap["updated_on"] = recordListItem.UpdatedOn` 赋值
- ✅ 验收测试断言 `data.tencentcloud_dnspod_record_list.foo.record_list.0.updated_on` 与 `data.tencentcloud_dnspod_record_list.foo.instance_list.0.updated_on` 至少有一条非空（使用 `TestMatchResourceAttr` 或 `TestCheckResourceAttrSet`）

---

#### Scenario: 单条记录 UpdatedOn 缺失时不影响其它字段

- **WHEN** SDK 返回的某条 `RecordListItem.UpdatedOn == nil`
- **THEN**
  - Read 跳过该条 `updated_on` 赋值
  - 该条目其它字段（`record_id`/`value`/`name` 等）正常填充
  - 整次 `terraform plan` / `apply` 流程不报错

**验收标准**：
- ✅ Read 中存在 `if recordListItem.UpdatedOn != nil { ... }` 守护
- ✅ 当云端返回部分条目无 `UpdatedOn` 时，结果列表长度与其它字段不受影响

---

### Requirement: REQ-DNSPOD-RECORD-UPDATED-ON-003 - 文档与示例同步

Documentation MUST be updated to reflect the new attribute on both the resource and the data source.

文档与示例必须同步更新，使用户能在网站文档中看到 `updated_on` 字段说明。

#### Scenario: 资源 markdown 与 website 文档展示 updated_on

- **WHEN** 用户在 Terraform Registry 网站查看 `tencentcloud_dnspod_record` 文档
- **THEN**
  - "Attributes Reference" 段落包含 `updated_on - Last update time of the record.`
  - `make doc` 已根据 `tencentcloud/services/dnspod/resource_tc_dnspod_record.md` 与 schema description 自动生成 `website/docs/r/dnspod_record.html.markdown`

**验收标准**：
- ✅ `website/docs/r/dnspod_record.html.markdown` 中可 grep 到 `updated_on`
- ✅ `make doc` 通过

---

#### Scenario: 数据源 markdown 与 website 文档展示嵌套 updated_on

- **WHEN** 用户在 Terraform Registry 网站查看 `tencentcloud_dnspod_record_list` 数据源文档
- **THEN**
  - `record_list` 与 `instance_list` 的 Attributes Reference 中均含 `updated_on - Update time.`
  - `website/docs/d/dnspod_record_list.html.markdown` 可 grep 到 `updated_on`

**验收标准**：
- ✅ 网站文档两处嵌套块均出现 `updated_on`
- ✅ `make doc` 通过
