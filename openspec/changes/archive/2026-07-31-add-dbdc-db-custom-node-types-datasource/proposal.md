## Why

用户需要通过 Terraform 查询 DB Custom（专用数据库）节点支持的机型列表信息，以便在创建节点时选择合适的机型规格。当前 Provider 已支持 DB Custom 集群、节点等数据源，但缺少查询可用节点机型信息的数据源，用户无法在 Terraform 配置中基于机型规格（CPU、内存、磁盘类型等）进行筛选和决策。

## What Changes

- 新增 Data Source: `tencentcloud_dbdc_db_custom_node_types`
- 实现对 DBDC API `DescribeDBCustomNodeTypes` 接口的调用（查询节点支持的机型列表）
- 支持通过 `filters` 进行过滤查询：
  - `filters`（可选）: 支持按 `region`、`zone`、`node-family`、`node-type` 过滤
  - `result_output_file`（可选）: 输出结果到文件
- 返回节点机型信息列表 `node_type_set`，每个机型包含：
  - `zone`: 可用区标识
  - `node_type`: 机型标识（如 DB.SA5.2XLARGE32）
  - `node_family`: 机型系列（如 DB.AT5、DB.SA5）
  - `cpu`: CPU 核数
  - `memory`: 内存大小（GiB）
  - `status`: 机型售卖状态（SELL / SOLD_OUT）
  - `system_disk_types`: 允许的系统盘类型列表
  - `data_disk_types`: 允许的数据盘类型列表

## Capabilities

### New Capabilities
- `dbdc-db-custom-node-types-datasource`: DB Custom 节点机型列表查询数据源，通过 `DescribeDBCustomNodeTypes` 接口查询可用节点机型规格及其售卖状态

### Modified Capabilities
<!-- 无需修改现有 capability 的 spec-level 行为 -->

## Impact

- **新增能力**: DB Custom 节点机型信息查询
- **受影响的服务**: dbdc (tencentcloud/services/dbdc)
- **新增文件**:
  - `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_types.go`
  - `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_types.md`
  - `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_types_test.go`
  - Provider 注册代码需在 DataSourcesMap 中添加此数据源
  - `service_tencentcloud_dbdc.go` 新增 `DescribeDBCustomNodeTypesByFilter` 方法
- **API 依赖**:
  - DBDC API v20201029: `DescribeDBCustomNodeTypes`（查询节点支持的机型列表，无分页）
- **兼容性**: 无破坏性变更，纯新增功能
