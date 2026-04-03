# 任务清单：add-region-datasources

## 1. 新建 service 层

**文件**: `tencentcloud/services/region/service_tencentcloud_region.go`

- [x] 声明 `package region`，定义 `RegionService` struct
- [x] 实现 `DescribeProductsByFilter`（带分页：for 循环 + 每次 retry，退出条件 len < limit）
- [x] 实现 `DescribeRegionsByFilter`（无分页，单次 retry）
- [x] 实现 `DescribeZonesByFilter`（无分页，单次 retry）
- [x] 执行 `go fmt ./tencentcloud/services/region/`

---

## 2. 新增 tencentcloud_products 数据源

- [x] `data_source_tc_products.go`：Schema + Read 函数
- [x] `data_source_tc_products.md`：文档示例
- [x] `data_source_tc_products_test.go`：测试

---

## 3. 新增 tencentcloud_regions 数据源

- [x] `data_source_tc_regions.go`：Schema（product Required，scene Optional，region_list 7字段）+ Read 函数
- [x] `data_source_tc_regions.md`：文档示例
- [x] `data_source_tc_regions_test.go`：测试（product="cvm"）

---

## 4. 新增 tencentcloud_zones 数据源

- [x] `data_source_tc_zones.go`：Schema（product Required，scene Optional，zone_list 8字段）+ Read 函数
- [x] `data_source_tc_zones.md`：文档示例
- [x] `data_source_tc_zones_test.go`：测试（product="cvm"）

---

## 5. 注册数据源

**文件**: `tencentcloud/provider.go`

- [x] import 新增 `regionpkg` 别名
- [x] 在 `DataSourcesMap` 中追加三个数据源注册

---

## 6. 编译验证

- [x] `go build ./tencentcloud/services/region/` 通过
- [x] `go build ./tencentcloud/` 通过
- [x] `go fmt ./tencentcloud/services/region/` 完成

---

## 总结

- **状态**: 🎉 所有任务已完成
- **风险等级**：低（纯新增）
- **破坏性变更**：无
