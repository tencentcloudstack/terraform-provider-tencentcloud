## Context

`tencentcloud_teo_zones` 数据源封装了 TEO（EdgeOne）`DescribeZones` 接口，用于查询当前 APPID 下授权的站点列表。现有 `zones` 输出字段已覆盖站点的加速区域、接入类型、标签、计费资源、NS 信息、归属权验证等，但未包含「版本管理配置组工作模式」信息。

腾讯云 TEO SDK（`teo/v20220901`）的 `Zone` 结构体已新增 `WorkModeInfos []*ConfigGroupWorkModeInfo` 字段，其中 `ConfigGroupWorkModeInfo` 包含：
- `ConfigGroupType *string`：配置组类型（`l7_acceleration` / `edge_functions` / `web_security`）
- `WorkMode *string`：工作模式（`immediate_effect` / `version_control`）

该字段描述站点各配置模块按配置组维度开启的工作模式（即时生效 vs 版本管理），参考 [版本管理](https://cloud.tencent.com/document/product/1552/113690)。

当前状态：
- `data_source_tc_teo_zones.go` 的 `zones` 嵌套 schema 不含 `work_mode_infos` 字段
- Read 函数 `dataSourceTencentCloudTeoZonesRead` 在遍历 `respData`（`[]*teov20220901.Zone`）时未读取 `WorkModeInfos`
- SDK 已 vendored，`Zone.WorkModeInfos` 与 `ConfigGroupWorkModeInfo` 结构体可直接使用，无需变更 vendor

约束：
- 必须保持向后兼容：仅新增 Computed 输出字段，不改变现有输入参数与已有输出字段
- 遵循现有 Read 函数的 nil 判断与字段映射模式（先判断切片非 nil 再遍历，子字段先判断非 nil 再 set）

## Goals / Non-Goals

**Goals:**
- 在 `zones` 嵌套 schema 中新增 `work_mode_infos`（TypeList，Computed）输出字段，每个元素含 `config_group_type`（TypeString，Computed）与 `work_mode`（TypeString，Computed）
- 在 Read 函数中将 `Zone.WorkModeInfos` 展平为 `work_mode_infos` 列表，遵循现有 nil 判断模式
- 同步更新 `.md` 示例文档
- 补充 gomonkey mock 单元测试覆盖 `WorkModeInfos` 读取

**Non-Goals:**
- 不新增或修改任何输入参数（filters/order/direction 保持不变）
- 不修改 `DescribeZones` API 调用逻辑或分页逻辑
- 不新增独立的 `work_mode` 资源或数据源
- 不改变数据源 ID 生成逻辑（仍基于 zone_id 列表 hash）

## Decisions

### Decision 1: `work_mode_infos` 作为 `zones` 下的嵌套 TypeList

**选择**：在 `zones` 的 `Elem` Resource schema 中新增 `work_mode_infos`（TypeList，Computed），`Elem` 为 Resource，内含 `config_group_type` 与 `work_mode` 两个 TypeString Computed 字段。

**理由**：
- `WorkModeInfos` 在 SDK 中是 `[]*ConfigGroupWorkModeInfo`（列表），一个站点可同时有多个配置组（七层加速、边缘函数、Web 防护）各自的工作模式
- 与现有 `zones` 内其他嵌套列表字段（如 `tags`、`resources`、`vanity_name_servers_ips`）的扁平模式一致
- 字段描述使用 SDK 中 `ConfigGroupWorkModeInfo` 的注释说明取值

### Decision 2: 展平逻辑遵循现有 nil 判断模式

**选择**：在 Read 函数遍历 `zones` 时，新增如下逻辑（与 `tags` / `resources` 等字段保持一致）：

```go
workModeInfosList := make([]map[string]interface{}, 0, len(zones.WorkModeInfos))
if zones.WorkModeInfos != nil {
    for _, workModeInfo := range zones.WorkModeInfos {
        workModeInfosMap := map[string]interface{}{}
        if workModeInfo.ConfigGroupType != nil {
            workModeInfosMap["config_group_type"] = workModeInfo.ConfigGroupType
        }
        if workModeInfo.WorkMode != nil {
            workModeInfosMap["work_mode"] = workModeInfo.WorkMode
        }
        workModeInfosList = append(workModeInfosList, workModeInfosMap)
    }
    zonesMap["work_mode_infos"] = workModeInfosList
}
```

**理由**：
- 与现有字段展平模式完全一致（先判断切片非 nil，遍历时子字段先判断非 nil 再 set）
- `WorkModeInfos` 为 `omitnil` 字段，可能返回 null，nil 判断避免写入空 map

### Decision 3: 单元测试使用 gomonkey mock

**选择**：在 `data_source_tc_teo_zones_test.go` 中新增测试用例，使用 gomonkey mock `DescribeTeoZonesByFilter` service 方法（或直接 mock 其调用的 API），构造含 `WorkModeInfos` 的 `[]*teov20220901.Zone` 响应，验证 `work_mode_infos` 被正确展平。

**理由**：
- 数据源为 RESOURCE_KIND_DATASOURCE，按项目要求新增字段时使用 gomonkey mock 进行业务逻辑单测
- 覆盖 `WorkModeInfos` 非空（含多配置组）与空两种场景

## Risks / Trade-offs

- **Risk**：`WorkModeInfos` 字段在部分历史 SDK 版本或老站点上可能返回 null/空 → **Mitigation**：Read 中先判断 `zones.WorkModeInfos != nil` 再遍历，空时不写 `work_mode_infos` key（与现有字段一致），不影响其他字段输出
- **Trade-off**：新增 Computed 输出字段后，存量 state 首次 refresh 会多出 `work_mode_infos` 字段 → 可接受，Computed 字段不产生 plan diff
