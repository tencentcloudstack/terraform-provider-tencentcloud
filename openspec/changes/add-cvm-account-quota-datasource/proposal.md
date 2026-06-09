## Why

用户需要查询腾讯云账户在不同地域和可用区的 CVM 配额详情，包括后付费、预付费、竞价实例配额，以及镜像配额和置放群组配额。这些信息对于资源规划、配额管理和自动化部署至关重要。目前 Provider 缺少对应的数据源支持。

## What Changes

- 新增数据源 `tencentcloud_cvm_account_quota`，调用 CVM DescribeAccountQuota API
- 支持通过可用区 (zone) 和配额类型 (quota-type) 进行过滤查询
- 返回 AppId 和完整的 AccountQuotaOverview 配额数据
- 包含后付费、预付费、竞价实例、镜像、置放群组等各类配额信息

## Capabilities

### New Capabilities
- `cvm-account-quota-datasource`: 查询 CVM 账户配额详情的数据源，支持按可用区和配额类型过滤

### Modified Capabilities
<!-- 无现有功能需要修改 -->

## Impact

**新增文件:**
- `tencentcloud/services/cvm/data_source_tc_cvm_account_quota.go` - 数据源实现
- `tencentcloud/services/cvm/data_source_tc_cvm_account_quota_test.go` - 测试文件
- `tencentcloud/services/cvm/data_source_tc_cvm_account_quota.md` - 文档模板
- `website/docs/d/cvm_account_quota.html.markdown` - 网站文档

**修改文件:**
- `tencentcloud/provider.go` - 注册新数据源
- `tencentcloud/provider.md` - 添加数据源到文档列表

**依赖:**
- 已有的 tencentcloud-sdk-go CVM 包
- 无需新增外部依赖
