## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoPrefetchOriginLimit` 方法，调用 DescribePrefetchOriginLimit 接口，通过 ZoneId 和 Filters（domain-name、area）查询匹配的限速配置

## 2. 资源 Schema 与 CRUD 实现

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_prefetch_origin_limit_config.go`，定义 ResourceTencentCloudTeoPrefetchOriginLimit 资源 schema（zone_id、domain_name、area、bandwidth、enabled）
- [x] 2.2 实现 resourceTencentCloudTeoPrefetchOriginLimitConfigCreate 方法，调用 ModifyPrefetchOriginLimit（Enabled=on），设置联合 ID（zone_id#domain_name#area）
- [x] 2.3 实现 resourceTencentCloudTeoPrefetchOriginLimitConfigRead 方法，调用 service 层 DescribeTeoPrefetchOriginLimit，匹配并设置字段（domain_name、area、bandwidth）
- [x] 2.4 实现 resourceTencentCloudTeoPrefetchOriginLimitConfigUpdate 方法，检测 bandwidth 和 enabled 变更，调用 ModifyPrefetchOriginLimit 更新
- [x] 2.5 实现 resourceTencentCloudTeoPrefetchOriginLimitConfigDelete 方法，调用 ModifyPrefetchOriginLimit（Enabled=off）删除配置

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 tencentcloud_teo_prefetch_origin_limit 资源
- [x] 3.2 在 `tencentcloud/provider.md` 中添加 tencentcloud_teo_prefetch_origin_limit 资源条目

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_prefetch_origin_limit_config_test.go`，使用 gomonkey mock 云 API，编写 Create/Read/Update/Delete 单元测试用例

## 5. 文档

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_prefetch_origin_limit_config.md`，包含一句话描述、Example Usage、Import 示例（说明联合 ID 格式）
