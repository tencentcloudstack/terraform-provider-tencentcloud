## Context

GA2（全球加速2.0）产品已有 `tencentcloud_ga2_endpoint_group` 资源在 `tencentcloud/services/ga2/` 目录下。现需新增数据源 `tencentcloud_ga2_accelerate_regions`，用于查询可选加速区域列表。

云API `DescribeAccelerateRegions` 无入参，返回 `AcceleratorRegionSet` 列表，包含地域名称、可用性、地域标识、地区名称、是否中国地域、支持的ISP类型、是否腾讯地域等字段。

## Goals / Non-Goals

**Goals:**
- 实现数据源 `tencentcloud_ga2_accelerate_regions`，调用 `DescribeAccelerateRegions` API 返回加速区域列表
- 遵循项目现有的数据源实现模式（参考 `tencentcloud_igtm_instance_list`）
- 提供完整的单元测试（使用 gomonkey mock）
- 注册到 provider 并提供文档

**Non-Goals:**
- 不实现过滤参数（API 本身无入参）
- 不实现分页（API 无分页参数）
- 不支持 `result_output_file`（API 无入参，数据量固定）

## Decisions

1. **数据源 Schema 设计**: 由于 `DescribeAccelerateRegions` 无入参，数据源 Schema 仅包含一个 Computed 的 `accelerator_region_set` 列表字段，以及可选的 `result_output_file` 字段用于输出结果到文件。

2. **直接在 Read 函数中调用 API**: 由于该接口简单（无入参、无分页），不需要在 `service_tencentcloud_ga2.go` 中新增服务层方法，直接在 Read 函数中通过 retry 调用 API 即可。参考同产品已有资源的模式。

3. **测试方式**: 使用 gomonkey 对云API进行 mock，不依赖真实环境，确保单元测试可在 CI 中运行。

4. **ID 设置**: 由于数据源无入参，使用 `helper.BuildToken()` 生成唯一 ID。

## Risks / Trade-offs

- [风险] API 返回数据量可能较大 → 由于该接口返回的是固定的区域列表，数据量有限，无需分页处理
- [风险] API 无入参导致无法过滤 → 用户可在 Terraform 中使用 `for_each` 或 `locals` 进行二次过滤
