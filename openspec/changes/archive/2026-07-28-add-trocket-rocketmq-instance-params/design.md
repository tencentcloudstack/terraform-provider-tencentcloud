## Context

`tencentcloud_trocket_rocketmq_instance` 是一个 RESOURCE_KIND_GENERAL 资源，管理腾讯云 TROcket（RocketMQ 5.x）实例的完整生命周期。当前 Schema 仅覆盖创建后付费实例所需的最小参数集。

云 API `CreateInstance`（trocket v20230308）已支持预付费相关入参与可用区/主题数配置，但 Terraform Schema 尚未暴露。本次变更需要将以下五个入参引入 Schema，覆盖 Create / Read / Update 链路：

| Schema 字段 | 云 API Create 字段 | 云 API Describe 字段 | 云 API Modify 字段 | 类型 |
|---|---|---|---|---|
| `pay_mode` | `request.PayMode` (*int64) | `response.PayMode` (*string) | 无 | int |
| `renew_flag` | `request.RenewFlag` (*int64) | `response.RenewFlag` (*int64) | 无 | int |
| `time_span` | `request.TimeSpan` (*int64) | 无对应字段 | 无 | int |
| `max_topic_num` | `request.MaxTopicNum` (*int64) | 无直接字段（有 TopicNumLimit/TopicNumUpperLimit） | `request.MaxTopicNum` (*int64) | int |
| `zone_ids` | `request.ZoneIds` ([]*int64) | `response.ZoneIds` ([]*int64) | `request.ZoneIds` ([]*string) | list[int] |

约束：必须保持向后兼容（纯新增 Optional 字段，不影响已有 state）；所有 SDK 调用使用 `tccommon.ReadRetryTimeout`/`WriteRetryTimeout` 并通过 `tccommon.RetryError` 包装错误。

## Goals / Non-Goals

**Goals:**
- 在 Schema 中新增 `pay_mode`、`renew_flag`、`time_span`、`max_topic_num`、`zone_ids` 五个 Optional 字段。
- Create 操作将五个字段填充到 `CreateInstanceRequest`。
- Read 操作从 `DescribeInstance` 响应回填 `pay_mode`、`renew_flag`、`zone_ids`。
- Update 操作支持 `max_topic_num`、`zone_ids` 原地更新；`pay_mode`/`renew_flag`/`time_span` 作为不可变参数。
- 补充单元测试与资源文档。

**Non-Goals:**
- 不修改既有任何 Schema 字段（instance_type、name、sku_code、remark、vpc_id、subnet_id、enable_public、bandwidth、ip_rules、message_retention 等）。
- 不新增 `_extension.go` 文件。
- 不处理 `time_span` 到期续费的定时任务类逻辑。

## Decisions

### Decision 1: `pay_mode`、`renew_flag`、`time_span` 标记为不可变（immutable）

Create 入参 `PayMode`、`RenewFlag`、`TimeSpan` 在 `ModifyInstanceRequest` 中不存在，云 API 不支持创建后修改付费模式/续费标志/购买时长。因此将这三个字段加入现有 Update 方法的 `immutableArgs` 数组，一旦变更返回 `fmt.Errorf("argument %s cannot be changed", v)`。

**备选方案**：将这些字段设为 `ForceNew`，但重建 RocketMQ 实例代价过高（数据丢失、实例 ID 变更），且付费模式变更本身不是合法运维操作，故选择 immutable 报错而非 ForceNew。

### Decision 2: `max_topic_num`、`zone_ids` 支持原地更新

`ModifyInstanceRequest` 同时包含 `MaxTopicNum`（*int64）和 `ZoneIds`（[]*string）。因此 `max_topic_num`、`zone_ids` 可在 Update 方法中原地修改，纳入 `request1`（`ModifyInstanceRequest`）的处理逻辑。

### Decision 3: `zone_ids` 类型差异处理

注意类型在 Create 与 Modify 接口间不一致：
- `CreateInstanceRequest.ZoneIds` 为 `[]*int64`
- `ModifyInstanceRequest.ZoneIds` 为 `[]*string`
- `DescribeInstanceResponseParams.ZoneIds` 为 `[]*int64`

Schema 统一使用 `TypeList` + `Elem: schema.TypeInt`。在 Create 中将 int 转为 `*int64`；在 Update 中将 int 转为 `*string`（`helper.IntToStr`）；在 Read 中将 `[]*int64` 转回 `[]interface{}` int。

### Decision 4: Read 回填 `pay_mode` 的枚举映射

云 API `DescribeInstance` 返回的 `PayMode` 为字符串枚举：`POSTPAID`（后付费）/ `PREPAID`（预付费），而 Create 入参 `PayMode` 为 int64（0=后付费，1=预付费）。为保持 Terraform Schema 类型一致（int），Read 时做映射：`POSTPAID` → 0，`PREPAID` → 1，其他值默认 0。

### Decision 5: Read 对 `time_span`、`max_topic_num` 的处理

`DescribeInstance` 响应无 `TimeSpan` 字段（购买时长不回显），故 `time_span` 在 Read 中不回填，仅在 state 中保持用户配置值（Optional 非 Computed，Terraform 不会因 read 未 set 而报 diff）。`max_topic_num` 无直接对应响应字段，同样不在 Read 中回填。

### Decision 6: Create 使用 `GetOk` 判断 Optional 字段

五个新字段均为 Optional。`pay_mode`、`renew_flag`、`time_span` 有默认值语义（0/0/1），但为避免与云 API 默认行为冲突并保证 0 值可显式传入，Create 中使用 `d.GetOk`（对 int 0 会判为 absent）。若用户未设置则不填入 request，由云 API 使用默认值。这与现有 `enable_public` 用 `GetOkExists` 区分 false/absent 的场景不同——此处无需区分 0 与 absent，因为 0 即等于云 API 默认值。

## Risks / Trade-offs

- **[Risk] `pay_mode` Read 枚举映射不覆盖未来新增计费模式** → 仅处理当前已知的 `POSTPAID`/`PREPAID`，未知值兜底为 0，并在注释中说明；云 API 若新增枚举后续再补。
- **[Risk] `time_span` 无法 Read 回填导致 state 与远端不同步** → 该字段为创建时一次性参数，生命周期内不变，不可变约束保证不会被误更新，state 中保留用户配置值即可。
- **[Trade-off] `zone_ids` Create/Modify 类型不一致增加代码复杂度** → 在 Create/Update 分别做 int→*int64 与 int→*string 转换，属于云 API 不规范的既定约束，无法规避。
