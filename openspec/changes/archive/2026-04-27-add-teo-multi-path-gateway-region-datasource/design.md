## Context

Terraform Provider for TencentCloud 当前缺少 TEO（边缘安全加速平台）多通道安全加速网关可用地域的数据源。用户在配置多通道安全加速网关时，需要先查询可用地域列表，但目前无法通过 Terraform 完成此操作。

云 API `DescribeMultiPathGatewayRegions`（teo/v20220901）已提供该查询能力，入参为 `ZoneId`（站点 ID），出参为 `GatewayRegions`（可用地域列表，每项包含 `RegionId`、`CNName`、`ENName`）。该接口为同步接口，无分页。

现有参考实现：`data_source_tc_teo_ip_region.go`，同样为 TEO 产品的简单查询类数据源，模式匹配度高。

## Goals / Non-Goals

**Goals:**
- 新增数据源 `tencentcloud_teo_multi_path_gateway_region`，支持按站点 ID 查询多通道安全加速网关可用地域列表
- 在 `provider.go` 和 `provider.md` 中注册新数据源
- 生成对应的 `.md` 文档文件
- 补充单元测试（使用 gomonkey mock）

**Non-Goals:**
- 不修改任何已有资源或数据源的 schema
- 不添加分页逻辑（API 本身不支持分页）
- 不支持异步轮询（API 为同步接口）

## Decisions

1. **数据源模式选择**: 采用 `data_source_tc_teo_ip_region.go` 的模式（简单查询、无分页），而非 `data_source_tc_teo_plans.go` 的分页模式。理由：`DescribeMultiPathGatewayRegions` 接口无分页参数，返回全量数据。

2. **Schema 设计**:
   - 入参：`zone_id`（Required, TypeString）—— 对应 API 请求的 `ZoneId`
   - 出参：`gateway_regions`（Computed, TypeList）—— 对应 API 响应的 `GatewayRegions`，嵌套 schema 包含 `region_id`、`cn_name`、`en_name`
   - 附加：`result_output_file`（Optional, TypeString）—— 遵循项目数据源惯例

3. **ID 生成策略**: 使用 `helper.DataResourceIdsHash()` 基于返回的地域 ID 列表生成数据源 ID，与 `data_source_tc_teo_ip_region.go` 保持一致。

4. **服务层方法**: 在 `service_tencentcloud_teo.go` 中新增 `DescribeTeoMultiPathGatewayRegionByFilter` 方法，封装 API 调用并添加 retry 处理。

5. **单元测试**: 使用 gomonkey mock 云 API 调用，仅测试业务逻辑，不使用 terraform 测试套件。

## Risks / Trade-offs

- [API 返回为空] → 当 `GatewayRegions` 为空时，数据源 ID 仍会被设置（基于空列表 hash），用户可通过 `gateway_regions` 长度为 0 判断无可用地域
- [ZoneId 无效] → 云 API 会返回错误，通过 `tccommon.RetryError()` 包装后返回给用户，terraform 会自然报错
