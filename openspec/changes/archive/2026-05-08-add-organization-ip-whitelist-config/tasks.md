# Tasks: Add Organization IP Whitelist Config Resource

## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/tco/` 目录下创建 `resource_tc_organization_ip_whitelist_config.go` 文件
- [x] 1.2 定义 Schema:
  - `zone_id`: Type=TypeString, Required=true, ForceNew=true, Description="Zone ID."
  - `ip_whitelist`: Type=TypeList(TypeString), Required=true, Description="IP whitelist entries."
- [x] 1.3 Schema 必须与 UpdateIPWhitelist API 入参保持一致

## 2. CRUD 函数实现

- [x] 2.1 Create 函数:
  - 读取 `zone_id` 字段
  - 使用 `zone_id` 作为资源 ID
  - 调用 Update 函数设置 IP 白名单
- [x] 2.2 Read 函数:
  - 调用 `GetIPWhitelist` API 获取当前 IP 白名单
  - 设置 `ip_whitelist` 到 state
  - 处理资源不存在的情况
- [x] 2.3 Update 函数:
  - 调用 `UpdateIPWhitelist` API 更新 IP 白名单
  - `zone_id` 不变，`ip_whitelist` 可更新（无 ForceNew）
- [x] 2.4 Delete 函数:
  - 调用 `UpdateIPWhitelist` API 清空 IP 白名单（设置为空列表）

## 3. Service 层

- [x] 3.1 在 `tencentcloud/services/tco/service_tencentcloud_organization.go` 中添加:
  - `DescribeOrganizationIPWhitelistConfigById` 方法
  - `DeleteOrganizationIPWhitelistConfigById` 方法

## 4. Provider 注册

- [x] 4.1 在 `tencentcloud/provider.go` 中注册新资源

## 5. 文档生成

- [x] 5.1 创建 `resource_tc_organization_ip_whitelist_config.md` 文档

## 6. 单元测试

- [x] 6.1 创建 `resource_tc_organization_ip_whitelist_config_test.go` 文件

## 7. 代码验证

- [x] 7.1 运行 `go build` 确保编译通过
- [x] 7.2 运行 `gofmt` 格式化代码