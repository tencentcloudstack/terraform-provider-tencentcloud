## Context

SQL Server 实例支持 SSL 加密功能，当前 Terraform Provider 中缺少对 SSL 配置的管理资源。已有类似 CONFIG 类型资源如 `tencentcloud_sqlserver_wan_ip_config` 和 `tencentcloud_postgresql_instance_ssl_config` 可作为参考。

云 API 接口：
- **读取**：`DescribeDBInstancesAttribute` — 查询实例附属属性，返回 `SSLConfig` 结构体包含 `Encryption`（SSL状态）、`SSLValidityPeriod`（证书有效期）、`SSLValidity`（证书有效性）、`IsKMS`（是否KMS加密）、`CMKId`（CMK密钥ID）、`CMKRegion`（CMK地域）
- **更新**：`ModifyDBInstanceSSL` — 开启/关闭/更新SSL加密，入参 `InstanceId`、`Type`（enable/disable/renew）、`WaitSwitch`（执行时机）、`IsKMS`（是否KMS加密）、`KeyId`、`KeyRegion`
- **异步轮询**：`ModifyDBInstanceSSL` 返回 `FlowId`，需通过 `DescribeFlowStatus` 轮询直到任务完成

已有服务层方法：`SqlserverService.DescribeSqlserverInstanceSslById` 已存在于 `service_tencentcloud_sqlserver.go`。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_sqlserver_db_instance_ssl_config` 资源，支持 SSL 开启、关闭、证书更新操作
- 支持 KMS 加密保护配置（IsKMS、KeyId、KeyRegion）
- 支持执行时机选择（立即执行或维护时间内执行）
- 正确处理异步操作，轮询 `DescribeFlowStatus` 等待 SSL 变更完成
- 支持 Import
- 符合 RESOURCE_KIND_CONFIG 模式：Create 调用 Update，Delete 为空操作

**Non-Goals:**
- 不管理 SSL 证书的下载和上传（云 API 不支持）
- 不管理 DescribeDBInstancesAttribute 返回的其他属性（如 TDEConfig、RegularBackup 等），这些由其他资源管理

## Decisions

### 1. 资源 ID 策略
使用 `instance_id` 作为资源 ID，因为 SSL 配置与实例一一对应。这与其他 CONFIG 类型资源（如 `postgresql_instance_ssl_config`、`sqlserver_wan_ip_config`）的模式一致。

### 2. Schema 设计
- `instance_id`：Required, ForceNew — 标识 SQL Server 实例
- `type`：Required — SSL 操作类型（enable/disable/renew）
- `wait_switch`：Optional, int — 执行时机（0-立即执行，1-维护时间内执行）
- `is_kms`：Optional, int — 是否 KMS 加密保护
- `key_id`：Optional, string — KMS CMK 密钥 ID（IsKMS=1 时必填）
- `key_region`：Optional, string — CMK 所属地域（IsKMS=1 时必填）
- Computed 字段：`encryption`、`ssl_validity_period`、`ssl_validity` — 从 Read 接口获取

### 3. 异步操作处理
`ModifyDBInstanceSSL` 返回 `FlowId`，需要：
1. 调用 `ModifyDBInstanceSSL` 获取 FlowId
2. 使用 `DescribeFlowStatus` 轮询，直到 Status == 0（SQLSERVER_TASK_SUCCESS）
3. 轮询超时使用 `tccommon.WriteRetryTimeout`

### 4. Read 方法
使用已有的 `SqlserverService.DescribeSqlserverInstanceSslById` 方法读取 SSL 配置。读取 `response.SSLConfig` 中的字段设置到 schema。当 SSLConfig 为 nil 或 Encryption 为 disable 时，资源仍然存在（CONFIG 类型资源不会"不存在"）。

### 5. Create/Delete 模式
- Create：设置 ID 为 instance_id，然后调用 Update
- Delete：空操作（CONFIG 类型资源，底层实例存在则配置存在）

## Risks / Trade-offs

- **[SSL 状态过渡期]** → SSL 开启/关闭/更新过程中有中间状态（enable_doing/disable_doing/renew_doing/wait_doing），Read 方法需要处理这些中间状态，不做资源删除处理
- **[异步操作延迟]** → 使用 WriteRetryTimeout 作为轮询超时时间，足够覆盖大部分场景
