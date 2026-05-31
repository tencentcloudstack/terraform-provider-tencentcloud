## Why

SQL Server 实例支持 SSL 加密功能，通过 `tencentcloud_sqlserver_db_instance_ssl_config` 资源以声明式方式管理 SSL 加密状态。

## What Changes

- `encryption` 参数从 Computed 改为 Required，表示 SSL 的期望状态（enable/disable）
- 移除 `type` 参数（操作型参数不适合声明式资源）
- Update 时根据 `encryption` 期望值决定 API 调用的 `Type`（enable→enable, disable→disable）
- 暂不支持 type=renew 操作
- 轮询逻辑从 DescribeFlowStatus 改为轮询 DescribeSqlserverInstanceSslById，直到 SSLConfig.Encryption 达到期望状态
- WaitSwitch 硬编码为 0（立即执行）
- Read 从 SSLConfig 读取 encryption（映射为 enable/disable）、is_kms、cmk_id、cmk_region

## Capabilities

### New Capabilities
- `sqlserver-db-instance-ssl-config`: 声明式管理 SQL Server 实例 SSL 加密状态

## Impact

- `tencentcloud/services/sqlserver/resource_tc_sqlserver_db_instance_ssl_config.go`
- `tencentcloud/services/sqlserver/resource_tc_sqlserver_db_instance_ssl_config_test.go`
- `tencentcloud/services/sqlserver/resource_tc_sqlserver_db_instance_ssl_config.md`
