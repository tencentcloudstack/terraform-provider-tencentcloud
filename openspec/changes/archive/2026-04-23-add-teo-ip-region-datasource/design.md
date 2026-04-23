## Context

Terraform Provider for TencentCloud 当前已支持 TEO 产品的多种数据源（zones、plans、origin_acl 等），但缺少 IP 归属查询数据源。TEO 的 `DescribeIPRegion` API 支持查询 IP 是否属于 EdgeOne 节点，这是用户在配置 CDN/安全策略时常用的判断依据。

当前状态：
- TEO 数据源文件位于 `tencentcloud/services/teo/`
- 服务层封装在 `TeoService` 结构体中（`service_tencentcloud_teo.go`）
- 数据源注册在 `provider.go` 中

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_ip_region` 数据源，支持通过 IP 列表查询归属信息
- 遵循现有 TEO 数据源的代码模式和组织结构
- 支持最大 100 个 IP 的批量查询
- 在 provider.go 和 provider.md 中正确注册数据源

**Non-Goals:**
- 不实现资源管理（CRUD），仅实现数据源（只读查询）
- 不修改现有 TEO 数据源的行为
- 不实现分页逻辑（该接口本身不支持分页）

## Decisions

### 1. Schema 设计
**决策**：入参 `ips` 使用 `schema.TypeList` + `schema.TypeString` 元素类型，出参 `ip_region_info` 使用 `schema.TypeList` + `schema.Resource` 元素类型。

**理由**：与云 API 的 `[]*string` 和 `[]*IPRegionInfo` 类型一致。参考现有 TEO 数据源（如 `tencentcloud_teo_zones`）的模式，`ip_region_info` 内嵌 `ip`（TypeString）和 `is_edge_one_ip`（TypeString）两个 Computed 字段。

### 2. 服务层方法
**决策**：在 `TeoService` 中新增 `DescribeTeoIPRegionByFilter` 方法，接受 `paramMap` 参数。

**理由**：遵循现有 TEO 数据源的服务层封装模式（如 `DescribeTeoZonesByFilter`），在服务层中构建请求、调用 SDK、返回结果。

### 3. ID 生成策略
**决策**：使用 `helper.DataResourceIdsHash(ids)` 生成数据源 ID，其中 ids 为查询结果中所有 IP 地址的列表。

**理由**：与现有 TEO 数据源保持一致，使用查询结果的标识生成唯一 ID。

### 4. 测试策略
**决策**：使用 gomonkey mock 云 API 进行单元测试，不使用 Terraform 测试套件。

**理由**：根据项目规范，新增资源使用 mock（gomonkey）方法进行单元测试，只测试业务代码逻辑。

## Risks / Trade-offs

- [API 限制] → `DescribeIPRegion` 单次最多查询 100 个 IP，在 schema 中通过 `MaxItems: 100` 约束
- [数据源无状态] → 数据源每次 Read 都会重新调用 API 查询，不做缓存，这是 Terraform 数据源的标准行为
