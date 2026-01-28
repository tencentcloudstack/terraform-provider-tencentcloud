# Proposal: Add tencentcloud_clickhouse_instances Data Source

## Change ID
`add-clickhouse-instances-datasource`

## Summary
新增 `tencentcloud_clickhouse_instances` 数据源，用于查询 TCHouse-C (ClickHouse) 实例列表。该数据源基于腾讯云 CDWCH 服务的 `DescribeInstancesNew` API 实现，支持多种过滤条件查询实例信息。

## Motivation
当前 Terraform Provider 中缺少一个专门用于批量查询 ClickHouse 实例的数据源。用户需要：
- **批量查询实例**：根据名称、ID、标签等条件查询多个 ClickHouse 实例
- **实例发现**：在 Terraform 配置中动态发现和引用已存在的实例
- **资源规划**：获取实例列表用于资源规划和管理
- **数据过滤**：支持分页和多种过滤条件以满足不同查询场景

目前只有 `resource_tc_clickhouse_instance` 资源可以管理单个实例，但没有对应的数据源用于查询和发现实例列表。

## Background

### 现有实现
- ✅ **Service 层已实现**：`CdwchService.DescribeInstancesNew()` 方法已存在于 `service_tencentcloud_cdwch.go`
- ✅ **SDK 支持**：`tencentcloud-sdk-go/tencentcloud/cdwch/v20200915` 已包含 `DescribeInstancesNew` API
- ✅ **Resource 已存在**：`resource_tc_clickhouse_instance` 资源已使用该 API
- ❌ **缺少 Data Source**：没有独立的数据源供用户查询实例列表

### API 信息
- **接口名称**：`DescribeInstancesNew`
- **接口文档**：https://cloud.tencent.com/document/product/1299/93328
- **API 版本**：`2020-09-15`
- **服务域名**：`cdwch.tencentcloudapi.com`
- **频率限制**：20 次/秒

### API 能力
**请求参数**：
- `SearchInstanceId` - 按实例 ID 搜索（支持精确匹配）
- `SearchInstanceName` - 按实例名称搜索（支持模糊匹配）
- `SearchTags` - 按标签搜索（支持多标签过滤）
- `Vips` - 按 VIP 地址搜索
- `Offset` / `Limit` - 分页参数
- `IsSimple` - 是否返回简化信息

**响应参数**：
- `TotalCount` - 实例总数
- `InstancesList` - 实例详情列表（包含 60+ 个字段）

## Proposed Changes

### 1. 创建 Data Source 文件
**文件**：`tencentcloud/services/cdwch/data_source_tc_clickhouse_instances.go`

**Schema 设计**：

**输入参数（过滤条件）**：
- `instance_id` - 实例 ID（可选，精确匹配）
- `instance_name` - 实例名称（可选，模糊匹配）
- `tags` - 标签过滤（可选，支持多标签）
- `vips` - VIP 地址列表（可选）
- `is_simple` - 是否返回简化信息（可选，默认 false）
- `result_output_file` - 结果输出文件路径（可选）

**输出参数（Computed）**：
- `instance_list` - 实例列表，包含以下核心字段：
  - **基础信息**：`instance_id`, `instance_name`, `status`, `status_desc`, `version`
  - **地域信息**：`region`, `zone`, `region_desc`, `zone_desc`
  - **网络信息**：`vpc_id`, `subnet_id`, `access_info`, `eip`, `ch_proxy_vip`
  - **计费信息**：`pay_mode`, `create_time`, `expire_time`, `renew_flag`
  - **配置信息**：`master_summary`, `common_summary` (节点规格汇总)
  - **高可用信息**：`ha`, `ha_zk`, `is_elastic`
  - **标签信息**：`tags`
  - **日志信息**：`has_cls_topic`, `cls_topic_id`, `cls_log_set_id`
  - **存储信息**：`cos_bucket_name`, `can_attach_cbs`, `can_attach_cos`
  - **组件信息**：`components`, `upgrade_versions`
  - **其他**：`kind`, `monitor`, `enable_xml_config`

### 2. 创建测试文件
**文件**：`tencentcloud/services/cdwch/data_source_tc_clickhouse_instances_test.go`

**测试场景**：
- 基础查询测试（不带过滤条件）
- 按实例 ID 过滤测试
- 按实例名称过滤测试
- 按标签过滤测试
- 分页测试

### 3. 创建文档文件
**文件**：`tencentcloud/services/cdwch/data_source_tc_clickhouse_instances.md`

**文档内容**：
- 使用示例（基础查询、带过滤条件查询）
- 参数说明
- 属性参考
- 注意事项

### 4. 注册 Data Source
**文件**：`tencentcloud/services/cdwch/extension_cdwch.go`

在 `GetResources()` 返回的 map 中添加：
```go
"tencentcloud_clickhouse_instances": DataSourceTencentCloudClickhouseInstances(),
```

### 5. Service 层（无需修改）
`CdwchService.DescribeInstancesNew()` 方法已存在，但当前实现仅支持按 `instance_id` 查询。需要在数据源中**直接调用 SDK API** 以支持更多过滤条件。

## Implementation Details

### Schema 映射策略
由于 API 返回字段非常多（60+ 字段），采用以下策略：
1. **核心字段映射**：实例基础信息、网络、计费、状态等常用字段全部映射
2. **嵌套对象处理**：对于复杂对象（如 `MasterSummary`, `CommonSummary`），映射为 TypeList
3. **可选字段**：部分字段可能为空，使用 `Computed: true` 标记
4. **类型转换**：API 返回的指针类型需要安全解引用

### 过滤逻辑
- **在 Terraform 层过滤**：对于 API 不支持的过滤条件（如 `result_output_file`），在数据源读取后进行二次过滤
- **在 API 层过滤**：充分利用 API 提供的过滤参数减少数据传输

### 分页处理
- 默认不设置 `Limit`，返回所有结果
- 如果未来需要支持，可在 schema 中添加 `limit` 和 `offset` 参数

## User Impact

### Benefits
- ✅ **批量查询能力**：用户可以一次性查询多个实例
- ✅ **灵活过滤**：支持按 ID、名称、标签、VIP 等多种条件过滤
- ✅ **动态引用**：在 Terraform 配置中动态发现和引用实例
- ✅ **完整信息**：返回实例的完整详细信息
- ✅ **符合规范**：遵循 Terraform Provider 的数据源设计模式

### 使用示例
```hcl
# 查询所有实例
data "tencentcloud_clickhouse_instances" "all" {}

# 按实例 ID 查询
data "tencentcloud_clickhouse_instances" "by_id" {
  instance_id = "cdwch-xxxxxx"
}

# 按实例名称查询
data "tencentcloud_clickhouse_instances" "by_name" {
  instance_name = "my-clickhouse"
}

# 按标签查询
data "tencentcloud_clickhouse_instances" "by_tag" {
  tags = {
    env = "production"
    app = "analytics"
  }
}

# 输出实例列表
output "instances" {
  value = data.tencentcloud_clickhouse_instances.all.instance_list
}
```

### Breaking Changes
**无** - 这是纯新增功能，不影响现有资源和数据源。

## Implementation Complexity
**中等** - 需要：
1. 创建数据源文件（约 300-400 行代码）
2. 定义完整的 Schema（60+ 字段映射）
3. 实现数据读取和转换逻辑
4. 编写测试用例（3-5 个场景）
5. 编写文档和使用示例

预计工作量：**1-2 天**

## Success Criteria
- [x] Data source 文件创建完成
- [x] Schema 定义完整，涵盖核心字段
- [x] 支持所有 API 提供的过滤条件
- [x] 数据读取逻辑正确，类型转换安全
- [x] 测试用例覆盖主要查询场景
- [x] 文档清晰，包含多个使用示例
- [x] 代码通过 `go fmt` 和 `make lint`
- [x] 验收测试通过
- [x] OpenSpec 验证通过

## Alternatives Considered

### Alternative 1: 扩展现有 data source
**考虑**：扩展 `data_source_tc_clickhouse_spec` 或其他现有数据源  
**拒绝原因**：
- `data_source_tc_clickhouse_spec` 用于查询规格信息，与实例查询是不同的领域
- 创建独立数据源更符合单一职责原则
- 避免单个数据源过于复杂

### Alternative 2: 最小化字段映射
**考虑**：只映射 10-15 个最核心字段  
**拒绝原因**：
- API 返回的大部分字段都有实际使用价值
- 用户可能需要完整信息做决策
- 后续扩展会造成破坏性变更

### Alternative 3: 使用现有 service 方法
**考虑**：直接使用 `CdwchService.DescribeInstancesNew()`  
**拒绝原因**：
- 现有方法仅支持按 `instance_id` 查询
- 需要修改 service 方法签名（可能影响现有调用）
- 在数据源中直接调用 SDK 更灵活

## Dependencies
- ✅ TencentCloud SDK Go (`cdwch/v20200915`) - 已在 vendor 中
- ✅ Service 层基础结构 - 已存在
- ✅ API 支持 - 已验证可用

## Related Changes
无。这是一个独立的功能添加。

## References
- 腾讯云 API 文档: https://cloud.tencent.com/document/product/1299/93328
- SDK 源码: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915/`
- 现有 Service 方法: `tencentcloud/services/cdwch/service_tencentcloud_cdwch.go:198`
- 相关资源: `tencentcloud/services/cdwch/resource_tc_clickhouse_instance.go`
- 参考数据源: `tencentcloud/services/cdwch/data_source_tc_clickhouse_spec.go`

## Timeline
- 第 1 天：实现数据源 Schema 定义和读取逻辑
- 第 2 天：编写测试用例、文档，代码审查和优化
