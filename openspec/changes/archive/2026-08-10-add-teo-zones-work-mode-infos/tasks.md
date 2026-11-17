## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/data_source_tc_teo_zones.go` 的 `zones` 嵌套 Resource schema 中，新增 `work_mode_infos` 字段（TypeList, Computed, Description: "版本管理配置组工作模式信息列表。"），其 `Elem` 为 Resource，内含 `config_group_type`（TypeString, Computed, Description 说明取值 `l7_acceleration`/`edge_functions`/`web_security`）与 `work_mode`（TypeString, Computed, Description 说明取值 `immediate_effect`/`version_control`）

## 2. Read 函数扩展

- [x] 2.1 在 `dataSourceTencentCloudTeoZonesRead` 遍历 `respData`（`[]*teov20220901.Zone`）的循环中，在现有字段展平逻辑之后（如 `ownership_verification` 之后、`zoneIds = append(...)` 之前），新增 `WorkModeInfos` 展平逻辑：先 `workModeInfosList := make([]map[string]interface{}, 0, len(zones.WorkModeInfos))`，再判断 `if zones.WorkModeInfos != nil`，遍历每个 `ConfigGroupWorkModeInfo`，按 nil 判断分别 set `config_group_type` 与 `work_mode`，最后 `zonesMap["work_mode_infos"] = workModeInfosList`
- [x] 2.2 确认展平逻辑与现有 `tags`/`resources`/`vanity_name_servers_ips` 等嵌套列表字段的 nil 判断与 map 组装模式完全一致

## 3. 文档同步

- [x] 3.1 在 `tencentcloud/services/teo/data_source_tc_teo_zones.md` 中更新说明，确保一句话描述中带上云产品名称 TEO（如 "Use this data source to query detailed information of TEO zones"），并在 Example Usage 中保留现有示例
- [x] 3.2 不要手动添加 `Argument Reference` 和 `Attribute Reference` 部分（由 `make doc` 工具自动生成）

## 4. 单元测试

- [x] 4.1 在 `tencentcloud/services/teo/data_source_tc_teo_zones_test.go` 中新增测试用例 `TestAccTencentCloudTeoZonesDataSource_workModeInfos`，使用 gomonkey mock `DescribeTeoZonesByFilter`（service 层方法）或其底层 API 调用，构造含 `WorkModeInfos`（两个元素：`l7_acceleration`/`immediate_effect` 与 `web_security`/`version_control`）的 `[]*teov20220901.Zone` 响应，验证 `work_mode_infos` 被正确展平
- [x] 4.2 新增测试用例覆盖 `WorkModeInfos` 为 nil 的场景，验证不报错且其他字段正常输出
- [x] 4.3 保证已有 `TestAccTencentCloudTeoZonesDataSource_basic` 测试用例不受影响（不修改其行为）

## 5. 验证

- [x] 5.1 代码正确性检查：确认新增的 `work_mode_infos` schema 字段与 `WorkModeInfos`（`[]*ConfigGroupWorkModeInfo`）的 `ConfigGroupType` / `WorkMode` 子字段在云 API `DescribeZones` 的 `Zone` 结构体中存在（vendor 路径 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`）
- [x] 5.2 检查 Read 函数中所有 error 返回已正确处理（本变更无新增 error 路径，仅字段展平）
