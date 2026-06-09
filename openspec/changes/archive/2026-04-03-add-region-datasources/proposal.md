# 变更提案：新增 Region 云产品三个 DataSource 资源

## 变更类型

**新功能** — 新增 `tencentcloud_products`、`tencentcloud_regions`、`tencentcloud_zones` 三个数据源，对应腾讯云地域管理系统（Region）API。

## Why

用户需要通过 Terraform 查询：
- 支持地域列表查询的云产品列表（`DescribeProducts`）
- 指定云产品支持的地域列表（`DescribeRegions`）
- 指定云产品支持的可用区列表（`DescribeZones`）

目前 Provider 中缺少这三个数据源，无法通过 IaC 方式查询地域相关信息。

## 接口信息

| DataSource | 接口名 | 文档 |
|------------|--------|------|
| `tencentcloud_products` | `DescribeProducts` | https://cloud.tencent.com/document/api/1596/77931 |
| `tencentcloud_regions` | `DescribeRegions` | https://cloud.tencent.com/document/api/1596/77930 |
| `tencentcloud_zones` | `DescribeZones` | https://cloud.tencent.com/document/api/1596/77929 |

SDK 包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region/v20220627`  
Client 方法：`meta.GetAPIV3Conn().UseRegionClient()`（已存在于 `connectivity/client.go`）

## 接口参数概要

### DescribeProducts
- 入参：`Limit`（int64, Optional）、`Offset`（int64, Optional）—— **有分页**
- 返回：`Products []*RegionProduct`（含 `Name *string`）、`TotalCount *int64`

### DescribeRegions
- 入参：`Product`（string, **Required**）、`Scene`（int64, Optional）—— **无分页**
- 返回：`RegionSet []*RegionInfo`（含 `Region`、`RegionName`、`RegionState`、`RegionTypeMC`、`LocationMC`、`RegionNameMC`、`RegionIdMC`）、`TotalCount *uint64`

### DescribeZones
- 入参：`Product`（string, **Required**）、`Scene`（int64, Optional）—— **无分页**
- 返回：`ZoneSet []*ZoneInfo`（含 `Zone`、`ZoneName`、`ZoneId`、`ZoneState`、`ParentZone`、`ParentZoneId`、`ParentZoneName`、`ZoneType`、`MachineRoomTypeMC`、`ZoneIdMC`）、`TotalCount *uint64`

## What Changes

### 新增文件

| 文件 | 说明 |
|------|------|
| `tencentcloud/services/region/data_source_tc_products.go` | products 数据源 |
| `tencentcloud/services/region/data_source_tc_products.md` | products 文档 |
| `tencentcloud/services/region/data_source_tc_products_test.go` | products 测试 |
| `tencentcloud/services/region/data_source_tc_regions.go` | regions 数据源 |
| `tencentcloud/services/region/data_source_tc_regions.md` | regions 文档 |
| `tencentcloud/services/region/data_source_tc_regions_test.go` | regions 测试 |
| `tencentcloud/services/region/data_source_tc_zones.go` | zones 数据源 |
| `tencentcloud/services/region/data_source_tc_zones.md` | zones 文档 |
| `tencentcloud/services/region/data_source_tc_zones_test.go` | zones 测试 |
| `tencentcloud/services/region/service_tencentcloud_region.go` | service 层封装 |

### 修改文件

| 文件 | 修改内容 |
|------|---------|
| `tencentcloud/provider.go` | 注册 3 个新数据源 |

### 向后兼容性

✅ 纯新增，不影响任何现有资源和数据源。
