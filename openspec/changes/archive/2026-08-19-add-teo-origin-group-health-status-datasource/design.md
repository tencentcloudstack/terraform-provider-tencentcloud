## Context

腾讯云 EdgeOne（TEO）提供了 `DescribeOriginGroupHealthStatus` API，用于查询负载均衡实例下源站组的健康状态。该 API 支持通过站点 ID（ZoneId）和负载均衡实例 ID（LBInstanceId）查询，可选传入源站组 ID 列表（OriginGroupIds）进行过滤。

当前 terraform-provider-tencentcloud 在 `tencentcloud/services/teo/` 目录下已有多个 TEO 相关的数据源实现（如 `data_source_tc_teo_origin_acl.go`、`data_source_tc_teo_default_certificate.go` 等），遵循统一的代码组织模式。新增数据源需要遵循现有的架构模式和编码规范。

云 API 响应结构如下：
- `OriginGroupHealthStatusList`（`[]*OriginGroupHealthStatusDetail`）：源站组健康状态列表
  - `OriginGroupId`（`*string`）：源站组 ID
  - `OriginHealthStatus`（`[]*OriginHealthStatus`）：综合决策的源站健康状态
    - `Origin`（`*string`）：源站
    - `Healthy`（`*string`）：健康状态（Healthy/Unhealthy/Undetected）
  - `CheckRegionHealthStatus`（`[]*CheckRegionHealthStatus`）：各健康检查区域的源站健康状态
    - `Region`（`*string`）：健康检查区域
    - `Healthy`（`*string`）：健康状态
    - `OriginHealthStatus`（`[]*OriginHealthStatus`）：源站健康状态

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_origin_group_health_status` 数据源，完整映射 `DescribeOriginGroupHealthStatus` API 的请求参数和响应字段
- 支持必填参数 ZoneId、LBInstanceId 和可选参数 OriginGroupIds
- 遵循 Provider 现有的代码组织结构和命名规范，参考 `data_source_tc_igtm_instance_list.go` 的实现模式
- 提供完整的单元测试（使用 gomonkey mock）和文档

**Non-Goals:**
- 不实现源站组健康状态的修改操作（该 API 为只读查询接口）
- 不对 API 返回的数据进行额外的聚合或统计
- 不实现分页（该 API 不支持分页参数）

## Decisions

### 1. Schema 设计
**决策**: 将 API 请求参数映射为数据源的输入参数，将 API 响应映射为 computed 字段

**理由**:
- `zone_id` 和 `lb_instance_id` 为 Required，因为云 API 中这两个参数为必填
- `origin_group_ids` 为 Optional，不填写时默认获取负载均衡下所有源站组的健康状态
- `origin_group_health_status_list` 为 Computed TypeList，包含嵌套的源站健康状态和检查区域健康状态

**Schema 映射**:
- 输入参数:
  - `zone_id`（Required, TypeString）→ `request.ZoneId`
  - `lb_instance_id`（Required, TypeString）→ `request.LBInstanceId`
  - `origin_group_ids`（Optional, TypeList of TypeString）→ `request.OriginGroupIds`
  - `result_output_file`（Optional, TypeString）→ 标准字段，用于导出结果
- 输出参数（Computed）:
  - `origin_group_health_status_list`（TypeList）→ `response.Response.OriginGroupHealthStatusList`
    - `origin_group_id`（TypeString）
    - `origin_health_status`（TypeList）
      - `origin`（TypeString）
      - `healthy`（TypeString）
    - `check_region_health_status`（TypeList）
      - `region`（TypeString）
      - `healthy`（TypeString）
      - `origin_health_status`（TypeList）
        - `origin`（TypeString）
        - `healthy`（TypeString）

### 2. 数据源 ID 生成
**决策**: 使用 `zone_id`、`lb_instance_id` 和 `origin_group_ids` 组合生成唯一 ID

**理由**:
- 数据源没有真实的资源 ID，需要生成唯一标识用于 Terraform state
- 使用参数组合的哈希值确保同一组查询参数对应同一个数据源实例

### 3. 重试机制
**决策**: 使用 `resource.Retry` 包装 API 调用，超时时间为 `tccommon.ReadRetryTimeout`

**理由**:
- 遵循项目规范，数据源 Read 方法需要处理最终一致性问题
- 在 retry 块内检查 API 返回是否为空，若为空返回 `NonRetryableError`，避免因云 API 短暂波动导致 state 中的 id 被清空

### 4. 单元测试策略
**决策**: 使用 gomonkey 对云 API 进行 mock，不使用 terraform 测试套件

**理由**:
- 遵循项目规范，新增 terraform 数据源的单元测试使用 mock 方法进行业务代码逻辑测试
- 避免依赖真实云资源，测试可独立运行

## Risks / Trade-offs

- **[API 返回空列表]** → 在 retry 块内检查 `response.Response.OriginGroupHealthStatusList` 是否为空，若为空返回 `NonRetryableError`，并在外层 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示，便于排障
- **[嵌套字段为 nil]** → 在设置每个字段前检查 nil，避免空指针错误
- **[参数变更未触发重新查询]** → `zone_id` 和 `lb_instance_id` 设为 Required，Terraform 会在参数变更时自动重新执行 Read
