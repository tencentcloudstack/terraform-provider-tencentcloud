## Why

`tencentcloud_teo_zones` 数据源当前仅返回站点的加速、接入、标签、计费资源等信息，但未暴露站点的「版本管理配置组工作模式」（`WorkModeInfos`）。EdgeOne 站点的各配置模块（七层加速、边缘函数、Web 防护）可按配置组维度开启「即时生效模式」或「版本管理模式」，用户在 Terraform 中无法读取该状态，导致无法在 IaC 流程中感知或校验站点的工作模式。腾讯云 TEO SDK 的 `DescribeZones` 接口已在 `Zone` 结构体中返回 `WorkModeInfos` 字段（类型为 `[]*ConfigGroupWorkModeInfo`，包含 `ConfigGroupType` 和 `WorkMode`），可直接在数据源 Read 中展平输出，无需新增 API 调用。

## What Changes

- 在 `tencentcloud_teo_zones` 数据源的 `zones` 嵌套 schema 中新增 `work_mode_infos` 列表字段（Computed，TypeList），用于输出站点各配置组的工作模式信息。
- `work_mode_infos` 每个元素包含两个 Computed 字段：
  - `config_group_type`（TypeString）：配置组类型，取值 `l7_acceleration`（七层加速配置组）、`edge_functions`（边缘函数配置组）、`web_security`（Web 防护配置组）。
  - `work_mode`（TypeString）：工作模式，取值 `immediate_effect`（即时生效模式）、`version_control`（版本管理模式）。
- 在 Read 函数中，将 `Zone.WorkModeInfos`（`[]*teov20220901.ConfigGroupWorkModeInfo`）展平为 `work_mode_infos` 列表，遵循现有 nil 判断与字段映射模式。
- 同步更新 `data_source_tc_teo_zones.md` 示例文档，补充新输出字段的说明。
- 在 `data_source_tc_teo_zones_test.go` 中补充单元测试，使用 gomonkey mock 云 API，覆盖 `WorkModeInfos` 字段的读取。

非破坏性：仅新增 Computed 输出字段，不改变现有输入参数与已有输出字段，不影响存量 TF 配置与 state。

## Capabilities

### New Capabilities
- `teo-zones-datasource`: 提供 `tencentcloud_teo_zones` 数据源的完整 schema 定义与 Read 行为规范，覆盖查询条件（filters/order/direction）、站点列表输出（含本次新增的 `work_mode_infos` 配置组工作模式信息）、分页处理、ID 生成与结果输出文件等需求。

### Modified Capabilities
<!-- 本次为该数据源首个 spec，不存在已有 spec 需修改 -->

## Impact

- **代码**：
  - `tencentcloud/services/teo/data_source_tc_teo_zones.go`（在 `zones` 嵌套 schema 中新增 `work_mode_infos` 字段；在 Read 函数中新增 `WorkModeInfos` 展平逻辑）
  - `tencentcloud/services/teo/data_source_tc_teo_zones_test.go`（新增覆盖 `work_mode_infos` 的单元测试用例）
  - `tencentcloud/services/teo/data_source_tc_teo_zones.md`（补充输出字段示例说明）
- **云 API**：复用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.Zone.WorkModeInfos`（`[]*ConfigGroupWorkModeInfo`，字段 `ConfigGroupType` / `WorkMode`），无需变更 vendor。
- **向后兼容**：纯新增 Computed 输出字段，不影响现有 schema 输入参数与已有输出，无 state drift。
- **文档**：需同步更新 `.md` 示例文件，由 `make doc` 流程生成 `website/docs/` 文档。
