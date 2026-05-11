# 任务清单：add-mysql-proxy-address-config

## 1. Service 层新增方法

**文件**: `tencentcloud/services/cdb/service_tencentcloud_mysql.go`

- [x] 新增 `DescribeMysqlProxyAddressConfig(ctx, instanceId, proxyGroupId, proxyAddressId string) (address *cdb.ProxyAddress, errRet error)` 方法
  - 复用 `DescribeMysqlProxyById` → 遍历 `ProxyAddress` 匹配 `ProxyAddressId`

---

## 2. 实现 Resource 主文件

**文件**: `tencentcloud/services/cdb/resource_tc_mysql_proxy_address_config.go`

- [x] 实现 `ResourceTencentCloudMysqlProxyAddressConfig()` schema 定义（Required + Optional 字段与 AdjustCdbProxyAddressRequest 对齐）
- [x] 实现 `resourceTencentCloudMysqlProxyAddressConfigCreate`（设置 ID = `instanceId#proxyGroupId`，调用 Update）
- [x] 实现 `resourceTencentCloudMysqlProxyAddressConfigRead`（调用 `DescribeMysqlProxyAddressConfig`，回填所有字段）
- [x] 实现 `resourceTencentCloudMysqlProxyAddressConfigUpdate`（先 `DescribeCdbProxyInfo` 校验唯一性，再调用 `AdjustCdbProxyAddress`）
- [x] 实现 `resourceTencentCloudMysqlProxyAddressConfigDelete`（no-op，返回 nil）

---

## 3. 生成 Resource 文档

**文件**: `tencentcloud/services/cdb/resource_tc_mysql_proxy_address_config.md`

- [x] 编写 Example Usage（完整 HCL 示例）
- [x] 编写 Import 说明（`instanceId#proxyGroupId`）

---

## 4. 生成单元测试

**文件**: `tencentcloud/services/cdb/resource_tc_mysql_proxy_address_config_test.go`

- [x] 实现 `TestAccTencentCloudMysqlProxyAddressConfigResource_basic`（Create check + Update check + ImportState）

---

## 5. 注册 Resource 到 Provider

**文件**: `tencentcloud/provider.go`

- [x] 注册 `"tencentcloud_mysql_proxy_address_config": cdb.ResourceTencentCloudMysqlProxyAddressConfig()`

---

## 6. 编译验证

- [x] 无编译 ERROR（linter 通过，仅 HINT 级别 deprecated 提示，符合 codebase 现有风格）

---

## 总结

- **预计工作量**：中等
- **风险等级**：低（Config 型资源，Delete 为 no-op）
- **破坏性变更**：无
- **状态**: 已完成
