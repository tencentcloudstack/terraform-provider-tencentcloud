# 设计文档：Region 三个 DataSource 资源

## 1. 文件结构

```
tencentcloud/services/region/           # 新建目录
├── service_tencentcloud_region.go      # service 层
├── data_source_tc_products.go          # products 数据源
├── data_source_tc_products.md
├── data_source_tc_products_test.go
├── data_source_tc_regions.go           # regions 数据源
├── data_source_tc_regions.md
├── data_source_tc_regions_test.go
├── data_source_tc_zones.go             # zones 数据源
├── data_source_tc_zones.md
└── data_source_tc_zones_test.go

tencentcloud/
└── provider.go                         # 修改：注册 3 个数据源
```

---

## 2. Service 层设计（service_tencentcloud_region.go）

### 包声明和 Service 结构

```go
package region

import (
    "context"
    "fmt"
    "log"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    regionv20220627 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region/v20220627"

    tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
    "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
    "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type RegionService struct {
    client *connectivity.TencentCloudClient
}
```

### DescribeProducts（有分页，for 循环 + retry）

```go
func (me *RegionService) DescribeProductsByFilter(ctx, param) ([]*regionv20220627.RegionProduct, error)
// 分页：offset/limit(100), for 循环, 每次循环内 resource.Retry
// 退出条件：len(Products) < limit
```

### DescribeRegions（无分页，单次 retry）

```go
func (me *RegionService) DescribeRegionsByFilter(ctx, param) ([]*regionv20220627.RegionInfo, error)
// 单次 resource.Retry
```

### DescribeZones（无分页，单次 retry）

```go
func (me *RegionService) DescribeZonesByFilter(ctx, param) ([]*regionv20220627.ZoneInfo, error)
// 单次 resource.Retry
```

---

## 3. DataSource Schema 设计

### data_source_tc_products.go

输入参数（均 Optional）：
- `result_output_file`（String）

输出字段：
- `product_list`（List，Computed）
  - `name`（String）：产品名称，如 `cvm`

### data_source_tc_regions.go

输入参数：
- `product`（String，**Required**）：待查询产品名称
- `scene`（Int，Optional）：场景控制，0 或 1
- `result_output_file`（String，Optional）

输出字段：
- `region_list`（List，Computed）
  - `region`（String）：地域标识，如 `ap-guangzhou`
  - `region_name`（String）：地域名称
  - `region_state`（String）：地域状态
  - `region_type_m_c`（Int）：控制台类型
  - `location_m_c`（String）：不同语言的地区
  - `region_name_m_c`（String）：控制台展示的地域描述
  - `region_id_m_c`（String）：控制台地域 ID

### data_source_tc_zones.go

输入参数：
- `product`（String，**Required**）
- `scene`（Int，Optional）
- `result_output_file`（String，Optional）

输出字段：
- `zone_list`（List，Computed）
  - `zone`（String）：可用区名称
  - `zone_name`（String）：可用区描述
  - `zone_id`（String）：可用区 ID
  - `zone_state`（String）：可用区状态
  - `parent_zone`（String）：父级 zone
  - `parent_zone_id`（String）：父级可用区 ID
  - `parent_zone_name`（String）：父级可用区描述
  - `zone_type`（String）：zone 类型

---

## 4. Read 函数设计（严格参照 igtm_instance_list 风格）

```go
func dataSourceTencentCloudProductsRead(d, meta) error {
    defer tccommon.LogElapsed("data_source.tencentcloud_products.read")()
    defer tccommon.InconsistentCheck(d, meta)()

    logId   = tccommon.GetLogId(nil)
    ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(...)
    service = RegionService{client: meta.GetAPIV3Conn()}

    paramMap := make(map[string]interface{})
    // 无过滤参数

    var respData []*regionv20220627.RegionProduct
    reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        result, e := service.DescribeProductsByFilter(ctx, paramMap)
        if e != nil { return tccommon.RetryError(e) }
        respData = result
        return nil
    })

    // 构建 product_list → d.Set
    // d.SetId(helper.BuildToken())
    // result_output_file 处理
}
```

regions 和 zones 的 Read 函数类似，区别：
- 从 `d.GetOk("product")` 和 `d.GetOk("scene")` 构建 paramMap

---

## 5. Service 层分页设计（DescribeProducts）

```go
var (
    offset int64 = 0
    limit  int64 = 100
)
for {
    request.Offset = &offset
    request.Limit = &limit
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        ratelimit.Check(request.GetAction())
        result, e := me.client.UseRegionClient().DescribeProducts(request)
        if e != nil { return tccommon.RetryError(e) }
        // log, nil check
        response = result
        return nil
    })
    if err != nil { errRet = err; return }

    ret = append(ret, response.Response.Products...)
    if len(response.Response.Products) < int(limit) { break }
    offset += limit
}
```

---

## 6. provider.go 注册位置

在 `DataSourcesMap` 末尾附近追加（与 igtm 相关数据源相邻）：

```go
"tencentcloud_products": regionpkg.DataSourceTencentCloudProducts(),
"tencentcloud_regions":  regionpkg.DataSourceTencentCloudRegions(),
"tencentcloud_zones":    regionpkg.DataSourceTencentCloudZones(),
```

import 别名：`regionpkg "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/region"`

---

## 7. 测试设计

每个数据源一个测试文件，测试函数：
- `TestAccTencentCloudProductsDataSource_basic`
- `TestAccTencentCloudRegionsDataSource_basic`（product="cvm"）
- `TestAccTencentCloudZonesDataSource_basic`（product="cvm"）
