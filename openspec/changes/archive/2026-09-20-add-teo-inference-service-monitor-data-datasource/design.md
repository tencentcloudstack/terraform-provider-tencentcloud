## Context

腾讯云 TEO（EdgeOne）提供 `DescribeInferenceServiceMonitorData` API 用于查询推理服务的监控数据，包括 CPU/GPU 使用率、实例数量、内存使用率等指标。该接口位于 SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`，vendor 中已支持。

当前 `tencentcloud/services/teo/` 目录下已有多个数据源实现，遵循统一的代码组织模式。本次新增的 `tencentcloud_teo_inference_service_monitor_data` 数据源与同包内的 `tencentcloud_teo_billing_data` 数据源性质相似（均为监控/时间序列数据查询，无分页字段），可作为主要实现参考。

云 API 参数与字段映射（已核对 vendor 中 models.go）：
- 入参（请求体 `DescribeInferenceServiceMonitorDataRequest`）：
  - `ZoneId`（*string）：站点 ID，必填
  - `ServiceIds`（[]*string）：推理服务 ID，最多 10 个，必填
  - `MetricNames`（[]*string）：指标列表，最多 10 个，必填
  - `StartTime`（*string）：开始时间，必填
  - `EndTime`（*string）：结束时间（EndTime - StartTime ≤ 30 天），必填
  - `Interval`（*string）：查询时间粒度（min/5min/hour/day），可选
- 出参（响应体 `DescribeInferenceServiceMonitorDataResponse.Response`）：
  - `TotalCount`（*int64）：结果总条数
  - `InferenceServiceMonitorRecords`（[]*InferenceServiceMonitorRecord）：监控记录列表
    - `ServiceId`（*string）：推理服务 ID
    - `MetricName`（*string）：指标名称
    - `InferenceServiceMonitorItems`（[]*InferenceServiceMonitorItem）：监控数据明细
      - `Timestamp`（*string）：时间点
      - `Value`（*float64）：具体数值

该接口为同步接口，无分页字段（不涉及 Offset/Limit），无需轮询。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_inference_service_monitor_data` 数据源，完整映射 `DescribeInferenceServiceMonitorData` API 的请求参数和响应字段
- 支持 ZoneId、ServiceIds、MetricNames、StartTime、EndTime 必填参数及 Interval 可选参数
- 返回监控记录列表，每条记录包含 ServiceId、MetricName 及监控数据明细（Timestamp、Value）
- 遵循 Provider 现有的代码组织结构、命名规范与错误处理模式
- 提供文档样例与单元测试（使用 mock 方式，不依赖真实云 API）

**Non-Goals:**
- 不实现推理服务监控数据的聚合或统计加工（仅原样返回 API 数据）
- 不暴露分页参数（该 API 本身无分页字段）
- 不实现推理服务的创建、修改、删除操作（这些属于资源管理范畴）

## Decisions

### 1. 实现模式：直接客户端调用（参考 data_source_tc_teo_billing_data.go）
**决策**: 在数据源 Read 函数中直接通过 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient()` 获取客户端并调用 `DescribeInferenceServiceMonitorDataWithContext`，不新增 service 层封装方法。

**理由**:
- 同包内的 `tencentcloud_teo_billing_data` 数据源性质相同（监控/时间序列数据、无分页），采用直接客户端调用模式
- 该 API 仅一个查询调用，无需复用，service 层封装会增加不必要的间接层
- 与现有 teo 数据源（billing_data、content_quota 等）保持一致

**备选方案**: 使用 service 层封装模式（参考 data_source_tc_teo_config_group_version_detail.go）
- **不采用理由**: service 层封装更适合会被多个资源/数据源复用的查询逻辑；本数据源查询逻辑单一，无需复用

### 2. Schema 设计
**决策**: 将 API 请求参数映射为数据源的输入参数（必填项 Required、可选项 Optional），将 API 响应映射为 computed 字段。监控记录列表使用 TypeList 嵌套两层 Resource（record → items）。

**Schema 映射**:
- 输入参数:
  - `zone_id`（TypeString, Required）← `ZoneId`
  - `service_ids`（TypeList of String, Required）← `ServiceIds`
  - `metric_names`（TypeList of String, Required）← `MetricNames`
  - `start_time`（TypeString, Required）← `StartTime`
  - `end_time`（TypeString, Required）← `EndTime`
  - `interval`（TypeString, Optional）← `Interval`
  - `result_output_file`（TypeString, Optional）：标准导出字段
- 输出参数（Computed）:
  - `inference_service_monitor_records`（TypeList）：监控记录列表，每条记录包含:
    - `service_id`（TypeString）← `InferenceServiceMonitorRecords.ServiceId`
    - `metric_name`（TypeString）← `InferenceServiceMonitorRecords.MetricName`
    - `inference_service_monitor_items`（TypeList）：监控数据明细，每项包含:
      - `timestamp`（TypeString）← `InferenceServiceMonitorItems.Timestamp`
      - `value`（TypeFloat）← `InferenceServiceMonitorItems.Value`

**理由**:
- 嵌套两层 TypeList 如实反映 API 的响应结构（record 列表 → items 列表）
- 列表型参数（service_ids、metric_names）使用 TypeList，与同包 billing_data 的 zone_ids 保持一致
- 不对列表再嵌套一层"总集合"schema（遵循扁平化输出要求），直接使用记录列表作为顶层输出

### 3. 重试与空响应处理
**决策**: 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用。在 retry 块内检查响应是否为空（`response == nil || response.Response == nil || len(response.Response.InferenceServiceMonitorRecords) == 0`），若为空则返回 `NonRetryableError`，不直接 `d.SetId("")`。

**理由**:
- 遵循数据源 Read 的标准重试模式与空响应保护要求
- 避免云 API 短暂波动导致本地 state 中 id 被清空造成数据丢失
- 让外层 retry 继续尝试，并以"重试耗尽"形式失败，便于人工介入

### 4. 资源 ID 生成
**决策**: 使用 `helper.BuildToken()` 生成唯一 ID（参考 data_source_tc_teo_igtm_instance_list 模式），或拼接 zone_id、start_time、end_time 等查询参数作为复合 ID。

**理由**:
- 数据源无持久化资源，ID 仅用于 Terraform state 标识
- 采用 `helper.BuildToken()` 生成唯一 token，简单可靠，与 igtm_instance_list 等数据源一致

### 5. 测试策略：mock 方式
**决策**: 单元测试使用 gomonkey 对云 API 客户端方法进行 mock，仅测试业务逻辑，不使用 Terraform 测试套件，不依赖真实云 API。

**理由**:
- 新增数据源按代码生成要求使用 mock 方式进行业务逻辑单元测试
- 不执行集成/端到端测试

## Risks / Trade-offs

**[风险] API 返回大量监控数据点导致 Terraform state 过大**
- → 缓解: 文档中建议用户合理设置时间范围（≤30 天）与时间粒度，避免一次查询过多数据点
- → 缓解: 指标与服务数量有 API 侧上限（各 10 个）

**[权衡] 直接客户端调用 vs service 层封装**
- 优势: 代码更紧凑，与同包监控类数据源一致，减少间接层
- 劣势: 查询逻辑不可被其他资源复用
- 决策: 本数据源查询逻辑单一且无复用需求，采用直接客户端调用

**[风险] 嵌套两层 TypeList 的状态管理复杂度**
- → 缓解: 严格按 nil 检查逐层 set，参考 billing_data 的 data 列表处理模式

## Migration Plan

**部署步骤**:
1. 实现数据源代码、文档样例与单元测试
2. 在 provider.go 注册数据源
3. 提交 PR，通过收尾阶段生成文档（make doc）

**回滚策略**:
- 纯新增功能，不影响现有资源与数据源，无需回滚

**文档更新**:
- 由收尾阶段通过 `make doc` 生成 website 文档
- 在 services/teo 下提供 data_source_tc_teo_inference_service_monitor_data.md 样例

## Open Questions

无待解决问题。API 字段已在 vendor 中核对完整，实现方案明确。