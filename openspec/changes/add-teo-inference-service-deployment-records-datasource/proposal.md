## Why

TEO（EdgeOne）推理服务（Inference Service）支持部署历史查询能力，用户可以通过云 API `DescribeInferenceServiceDeploymentRecords` 查询某个推理服务的部署历史列表，获取每次部署的操作类型、状态、耗时、配置快照以及是否为当前生效配置。当前 Terraform Provider 尚未提供对应的数据源，用户无法在基础设施即代码工作流中读取推理服务部署历史信息，难以对部署记录进行审计与引用。

## What Changes

- 新增数据源 `tencentcloud_teo_inference_service_deployment_records`，调用 TEO API `DescribeInferenceServiceDeploymentRecords` 查询推理服务部署历史列表
- 数据源入参：
  - `zone_id`（Required）：站点 ID
  - `service_id`（Required）：推理服务 ID
  - `sort_by`（Optional）：排序字段，取值 `create-time`，默认 `create-time`
  - `sort_order`（Optional）：排序方式，取值 `asc` / `desc`，默认 `desc`
  - `result_output_file`（Optional）：用于保存查询结果
- 数据源出参 `record_set`（TypeList, Computed）为部署历史记录列表，每个元素展开为平铺字段：
  - 部署记录信息：`record_id`、`operation`、`status`、`duration`、`create_time`、`active_status`
  - 推理服务配置 `inference_service_config`：`listen_port`、`request_paths`、`containers`、`resource_config`、`affinity_config`
    - `containers` 下：`image_type`、`tcr_repository_config`（`tcr_type`、`image`、`registry_id`、`region_name`）、`startup_command`、`environment_variables`（`key`、`value`）
    - `resource_config` 下：`scaling_mode`、`hardware_spec`、`hardware_spec_id`、`hardware_config`（`gpu_num`、`cpu_num`、`mem_size`、`disk_size`）、`auto_scaling_config`（`min_instance_count`、`scaling_policies`（`policy_name`、`policy_type`、`scheduled_scaling_policy`））、`manual_instance_config`（`fixed_instance_count`）、`concurrency`
    - `affinity_config` 下：`switch`、`affinity_mode`、`session_id_affinity_config`（`source`、`header_name`）
- 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增服务方法 `DescribeTeoInferenceServiceDeploymentRecordsByFilter`，内部实现分页（Limit 取云 API 最大值 100）自动获取全部数据
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新数据源
- 新增对应文档 `.md` 文件与单元测试文件 `*_test.go`（使用 gomonkey mock 方式测试业务逻辑）

## Capabilities

### New Capabilities
- `teo-inference-service-deployment-records-datasource`: 通过 `DescribeInferenceServiceDeploymentRecords` API 查询 TEO 推理服务部署历史列表的数据源

### Modified Capabilities

## Impact

- 新增文件：
  - `tencentcloud/services/teo/data_source_tc_teo_inference_service_deployment_records.go`
  - `tencentcloud/services/teo/data_source_tc_teo_inference_service_deployment_records_test.go`
  - `tencentcloud/services/teo/data_source_tc_teo_inference_service_deployment_records.md`
- 修改文件：
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`（新增服务方法）
  - `tencentcloud/provider.go`（注册数据源）
  - `tencentcloud/provider.md`（新增数据源声明）
- API 依赖：`DescribeInferenceServiceDeploymentRecords`，来自 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
- 兼容性：纯新增功能，无破坏性变更