# Change: Add TKE Addon Values Data Source

## Why

用户在使用 `tencentcloud_kubernetes_addon_config` 等资源安装 TKE 组件时，需要知道该组件当前生效的参数（`values`）以及该组件 chart 支持的默认参数（`default_values`）。目前 Provider 已经通过 `DescribeAddonValues` 接口在服务层封装了 `DescribeAddonValuesList` 方法（`service_tencentcloud_tke_addons.go`），但**没有暴露任何 Terraform data source**，用户只能在控制台或自行调用云 API 才能获取这些参数，无法在 Terraform 配置中以引用（`data.xxx.values`）的方式消费。

因此新增 data source `tencentcloud_tke_addon_values`，把 TKE `DescribeAddonValues` 接口的返回结果以 Terraform 数据源形式暴露。

## What Changes

- 新增 Data Source：`tencentcloud_tke_addon_values`
- 调用 TKE `DescribeAddonValues`（`tke/v20180525`）接口，复用已存在的服务层方法 `TkeService.DescribeAddonValuesList`
- 入参（Required）：
  - `cluster_id`：TKE 集群 ID
  - `addon_name`：addon（组件）名称
- 出参（Computed）：
  - `values`：addon 参数 JSON 字符串（已安装时返回按已设置参数 + chart 默认参数渲染后的结果；未安装时云 API 返回空值）
  - `default_values`：addon 支持的参数列表 JSON 字符串（即 chart 默认值）
- 通用出参：`result_output_file`（Optional，用于把结果写入文件）
- 新增文件：
  - `tencentcloud/services/tke/data_source_tc_tke_addon_values.go`
  - `tencentcloud/services/tke/data_source_tc_tke_addon_values.md`
  - `tencentcloud/services/tke/data_source_tc_tke_addon_values_test.go`（gomonkey mock 云 API 的单元测试）
- 在 `tencentcloud/provider.go` 的 `DataSourcesMap` 与 `tencentcloud/provider.md` 中注册 `tencentcloud_tke_addon_values`

## Capabilities

### New Capabilities

- `tke-addon-values-datasource`: 通过 Terraform data source 查询 TKE 指定集群下指定 addon 的当前参数（`values`）与 chart 默认参数（`default_values`）

### Modified Capabilities

无。本次仅新增 data source，不改变任何已有资源/数据源的行为与 schema。

## Impact

- **新增能力**：TKE addon 参数查询
- **受影响的服务**：TKE (`tencentcloud/services/tke`)
- **云 API 依赖**：`tke` `v20180525` 包中的 `DescribeAddonValues` 接口（已在 vendor 中存在，无需升级 SDK）
- **服务层依赖**：复用现有 `TkeService.DescribeAddonValuesList`（已存在于 `service_tencentcloud_tke_addons.go`），本次不修改其签名
- **注册点**：`tencentcloud/provider.go`、`tencentcloud/provider.md`
- **文档**：`website/docs/d/tke_addon_values.html.markdown` 由收尾阶段 `make doc` 生成，不手工创建
- **兼容性**：纯新增，无破坏性变更；不改动已有资源 schema 与 state
- **API 文档**：https://cloud.tencent.com/document/product/457/39048
