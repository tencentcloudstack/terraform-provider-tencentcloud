## ADDED Requirements

### Requirement: Data source schema for inference service deployment records
The data source `tencentcloud_teo_inference_service_deployment_records` SHALL define the following schema:

**Input parameters:**
- `zone_id` (TypeString, Required): 站点 ID
- `service_id` (TypeString, Required): 推理服务 ID
- `sort_by` (TypeString, Optional): 排序字段，取值 `create-time`，默认 `create-time`
- `sort_order` (TypeString, Optional): 排序方式，取值 `asc` / `desc`，默认 `desc`
- `result_output_file` (TypeString, Optional): 用于保存结果

**Computed output:**
- `record_set` (TypeList, Computed): 推理服务部署历史记录列表，每个元素平铺以下字段：
  - `record_id` (TypeString, Computed): 部署记录 ID
  - `operation` (TypeString, Computed): 部署操作类型（create/update/resume/stop）
  - `status` (TypeString, Computed): 部署状态（processing/succeeded/failed）
  - `duration` (TypeInt, Computed): 部署时长，单位秒
  - `create_time` (TypeString, Computed): 部署发起时间
  - `active_status` (TypeString, Computed): 是否当前生效配置（active/inactive）
  - `inference_service_config` (TypeList, Computed): 推理服务部署配置，包含：
    - `listen_port` (TypeInt, Computed): 监听端口
    - `request_paths` (TypeList, Computed, Elem TypeString): 请求路径列表
    - `containers` (TypeList, Computed): 容器配置列表，每个元素包含：
      - `image_type` (TypeString, Computed): 镜像类型
      - `tcr_repository_config` (TypeList, Computed): TCR 镜像仓库配置，包含 `tcr_type`、`image`、`registry_id`、`region_name`
      - `startup_command` (TypeString, Computed): 容器启动命令
      - `environment_variables` (TypeList, Computed): 环境变量列表，每个元素包含 `key`、`value`
    - `resource_config` (TypeList, Computed): 资源配置，包含：
      - `scaling_mode` (TypeString, Computed): 扩缩容方式（Auto/Manual）
      - `hardware_spec` (TypeString, Computed): 硬件规格标识（已废弃）
      - `hardware_spec_id` (TypeString, Computed): 硬件规格唯一标识 ID
      - `hardware_config` (TypeList, Computed): 硬件配置，包含 `gpu_num`、`cpu_num`、`mem_size`、`disk_size`
      - `auto_scaling_config` (TypeList, Computed): 自动伸缩配置，包含 `min_instance_count`、`scaling_policies`
      - `scaling_policies` (TypeList, Computed): 伸缩策略列表，每个元素包含 `policy_name`、`policy_type`、`scheduled_scaling_policy`
      - `manual_instance_config` (TypeList, Computed): 人工实例配置，包含 `fixed_instance_count`
      - `concurrency` (TypeInt, Computed): 单实例并发数
    - `affinity_config` (TypeList, Computed): 亲和性配置，包含：
      - `switch` (TypeString, Computed): 亲和总开关
      - `affinity_mode` (TypeString, Computed): 亲和方式
      - `session_id_affinity_config` (TypeList, Computed): 会话 ID 亲和配置，包含 `source`、`header_name`

#### Scenario: Data source with required parameters
- **WHEN** a user declares the data source with `zone_id` and `service_id`
- **THEN** Terraform SHALL call `DescribeInferenceServiceDeploymentRecords` API with the provided parameters and populate the `record_set` computed output

#### Scenario: Data source with optional sort parameters
- **WHEN** a user declares the data source with `sort_by` and `sort_order`
- **THEN** Terraform SHALL pass these as `SortBy` and `SortOrder` request parameters to the API

#### Scenario: Data source with result output file
- **WHEN** a user provides `result_output_file` parameter
- **THEN** the data source SHALL write the query results to the specified file path

### Requirement: Read function calls DescribeInferenceServiceDeploymentRecords API
The Read function SHALL call the `DescribeInferenceServiceDeploymentRecords` API via a service method from the `teo/v20220901` SDK package with `ZoneId` and `ServiceId` (Required) and optionally `SortBy` and `SortOrder` as request parameters.

#### Scenario: Successful API call
- **WHEN** the Read function is invoked with valid `zone_id` and `service_id`
- **THEN** it SHALL build a `DescribeInferenceServiceDeploymentRecordsRequest` with the provided parameters and call the API with `tccommon.ReadRetryTimeout` retry logic

#### Scenario: API call failure
- **WHEN** the API call fails with a retryable error
- **THEN** the Read function SHALL return a retry error using `tccommon.RetryError()`

### Requirement: Empty response handling for datasource
The Read function SHALL handle empty API responses safely within the retry block. When the service method returns an empty result (`nil` or zero-length `record_set`), the Read function SHALL NOT call `d.SetId("")` to clear state; instead it SHALL return a `NonRetryableError` so the outer retry continues, and SHALL log `[DATASOURCE] read empty, skip SetId` on the failure path.

#### Scenario: API returns empty record set
- **WHEN** the service method returns `nil` or an empty `record_set`
- **THEN** the Read function SHALL return a `NonRetryableError` and SHALL NOT clear the data source ID

### Requirement: Nil safety for nested response fields
The Read function SHALL check nil at each level of the nested response before accessing child fields. This includes checking `InferenceServiceConfig`, each `Containers` element, `TcrRepositoryConfig`, `EnvironmentVariables`, `ResourceConfig`, `HardwareConfig`, `AutoScalingConfig`, each `ScalingPolicies` element, `ScheduledScalingPolicy`, `ManualInstanceConfig`, `AffinityConfig`, and `SessionIdAffinityConfig` before reading their fields.

#### Scenario: Response with nil inference service config
- **WHEN** the API returns a record where `InferenceServiceConfig` is nil
- **THEN** the Read function SHALL NOT set `inference_service_config` and SHALL NOT panic

#### Scenario: Response with nil container tcr repository config
- **WHEN** a container element has nil `TcrRepositoryConfig`
- **THEN** the Read function SHALL NOT set `tcr_repository_config` for that container and SHALL NOT panic

#### Scenario: Response with nil resource config sub-objects
- **WHEN** `ResourceConfig` is present but `HardwareConfig` or `AutoScalingConfig` or `ManualInstanceConfig` is nil
- **THEN** the Read function SHALL NOT set the corresponding nested block and SHALL NOT panic

#### Scenario: Response with nil affinity config sub-objects
- **WHEN** `AffinityConfig` is present but `SessionIdAffinityConfig` is nil
- **THEN** the Read function SHALL NOT set `session_id_affinity_config` and SHALL NOT panic

### Requirement: Data source ID uses composite key
The data source SHALL set its ID as a composite of `zone_id` and `service_id` joined by `tccommon.FILED_SP`.

#### Scenario: Setting data source ID
- **WHEN** the Read function successfully retrieves data
- **THEN** the data source ID SHALL be set to `zone_id + FILED_SP + service_id`

### Requirement: Service method with automatic pagination
A service method `DescribeTeoInferenceServiceDeploymentRecordsByFilter` SHALL be added to `service_tencentcloud_teo.go` that accepts a `paramMap` with `ZoneId`, `ServiceId`, `SortBy`, and `SortOrder`, constructs the SDK request, calls the API with internal pagination (`Limit` fixed to the API maximum value 100 and `Offset` incremented until all records are fetched), merges all `RecordSet` results, and returns the combined `[]*teov20220901.InferenceServiceDeploymentRecord`. The pagination parameters SHALL NOT be exposed to the user.

#### Scenario: Service method fetches all pages
- **WHEN** the service method is called with `ZoneId` and `ServiceId`
- **THEN** it SHALL loop with `Offset` starting at 0 and `Limit` 100, appending each page's `RecordSet` until a page returns fewer than 100 records, and return the merged list

#### Scenario: Service method with sort parameters
- **WHEN** the service method is called with `SortBy` and `SortOrder` in the paramMap
- **THEN** it SHALL set these on the request for every paginated call

### Requirement: Provider registration
The data source SHALL be registered in `provider.go` under the dataSources map with key `tencentcloud_teo_inference_service_deployment_records` and value `teo.DataSourceTencentCloudTeoInferenceServiceDeploymentRecords()`.

#### Scenario: Data source registered in provider
- **WHEN** the provider is initialized
- **THEN** the dataSources map SHALL contain the entry `tencentcloud_teo_inference_service_deployment_records` mapped to `teo.DataSourceTencentCloudTeoInferenceServiceDeploymentRecords()`

### Requirement: Documentation file
A documentation file `data_source_tc_teo_inference_service_deployment_records.md` SHALL be created following the TEO doc pattern: a one-line description mentioning the TEO product name, an Example Usage section, and no `Argument Reference` / `Attribute Reference` sections (auto-generated). DATASOURCE resources SHALL NOT include an Import section.

#### Scenario: Documentation follows format rules
- **WHEN** the documentation file is created
- **THEN** it SHALL contain a one-line description with the TEO product name, an Example Usage section, and SHALL NOT contain Argument Reference or Attribute Reference sections

### Requirement: Unit tests using gomonkey mock
A unit test file `data_source_tc_teo_inference_service_deployment_records_test.go` SHALL be created using the gomonkey mock approach to mock the cloud API and test the Read function business logic. It SHALL NOT use the Terraform test suite.

#### Scenario: Unit tests cover Read function logic
- **WHEN** the unit tests are executed
- **THEN** they SHALL mock the service method / API call and verify the Read function correctly maps the nested response fields into the schema