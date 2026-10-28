## Why

用户需要通过 Terraform 查询 DLC（数据湖计算中心）TCLake 元数据实例的开通状态。当前 Provider 缺少对 `DescribeTCLakeMetaInstance` 接口的支持，该接口用于查询 TCLake 的开通状态（如 `Running` 表示开通成功）。通过 Data Source 的方式暴露该查询能力，可让用户在基础设施即代码场景中感知 TCLake 的开通进度与状态。

## What Changes

- 新增 Data Source: `tencentcloud_dlc_tc_lake_meta_instance`
- 实现对 DLC API `DescribeTCLakeMetaInstance` 接口的调用（包名: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125`）
- 该接口无入参，返回 `response.Response.Status`（开通状态枚举值，如 `Running`）
- 将 `status` 作为 computed 字段暴露给用户
- 支持 `result_output_file` 可选参数用于保存查询结果

## Capabilities

### New Capabilities
- `dlc-tc-lake-meta-instance-datasource`: 查询 DLC TCLake 元数据实例开通状态的 Data Source，调用 `DescribeTCLakeMetaInstance` 接口并返回 `status` 字段

### Modified Capabilities
<!-- 无需修改现有能力 -->

## Impact

- **新增能力**: DLC TCLake 开通状态查询 Data Source
- **受影响的服务**: DLC (tencentcloud/services/dlc)
- **新增文件**:
  - `tencentcloud/services/dlc/data_source_tc_dlc_tc_lake_meta_instance.go`
  - `tencentcloud/services/dlc/data_source_tc_dlc_tc_lake_meta_instance_test.go`
  - `tencentcloud/services/dlc/data_source_tc_dlc_tc_lake_meta_instance.md`
  - Provider 注册代码需要添加此 data source
- **修改文件**:
  - `tencentcloud/provider.go` 在 DataSourcesMap 中注册
  - `tencentcloud/provider.md` 追加资源名
  - `tencentcloud/services/dlc/service_tencentcloud_dlc.go` 添加服务层方法封装
- **API 依赖**:
  - DLC API v20210125: `DescribeTCLakeMetaInstance`（查询TCLake开通状态，无入参，出参 `Response.Status`）
