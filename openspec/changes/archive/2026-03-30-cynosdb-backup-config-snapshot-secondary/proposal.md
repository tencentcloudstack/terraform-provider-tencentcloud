# Proposal: Add snapshot_secondary_backup_config to tencentcloud_cynosdb_backup_config

## Overview

为 `tencentcloud_cynosdb_backup_config` 资源新增 `snapshot_secondary_backup_config` 字段，对应腾讯云 TDSQL-C MySQL API 中的二级快照备份配置参数（`SnapshotSecondaryBackupConfig`）。

## Motivation

### 当前问题

`tencentcloud_cynosdb_backup_config` 资源目前支持基础全量备份配置（`backup_time_beg`/`backup_time_end`/`reserve_duration`）和逻辑备份配置（`logic_backup_config`），但缺少对二级快照备份（Secondary Snapshot Backup）的支持。

二级快照备份是 TDSQL-C MySQL 的重要功能，允许用户配置独立于主备份的第二套快照备份策略，包括：
- 独立的备份时间窗口
- 独立的保留周期
- 跨地域备份支持
- 自定义备份触发策略

### 涉及的 API

- **修改接口**：[ModifyBackupConfig](https://cloud.tencent.com/document/api/1003/48090) — 请求参数 `SnapshotSecondaryBackupConfig`，类型 `SnapshotBackupConfig`
- **查询接口**：[DescribeBackupConfig](https://cloud.tencent.com/document/api/1003/48094) — 返回参数 `SnapshotSecondaryBackupConfig`，类型 `BackupConfigInfo`

> 注意：修改入参和查询返回使用了不同的数据结构：
> - `SnapshotBackupConfig`（ModifyBackupConfig 入参）：不含跨地域字段
> - `BackupConfigInfo`（DescribeBackupConfig 返回）：含 `CrossRegions`、`CrossRegionsEnable` 字段

## Proposed Solution

### Schema 新增字段

在现有 schema 中添加 `snapshot_secondary_backup_config`（TypeList，MaxItems:1），包含以下子字段：

| TF 字段名 | API 字段 | 类型 | 说明 |
|---|---|---|---|
| `backup_custom_auto_time` | `BackupCustomAutoTime` | Bool | 系统自动时间 |
| `backup_time_beg` | `BackupTimeBeg` | Int | 备份开始时间（秒） |
| `backup_time_end` | `BackupTimeEnd` | Int | 备份结束时间（秒） |
| `backup_week_days` | `BackupWeekDays` | List(String) | 备份星期数组 |
| `backup_interval_time` | `BackupIntervalTime` | Int | 备份间隔时间 |
| `reserve_duration` | `ReserveDuration` | Int | 保留时长（秒） |
| `backup_trigger_strategy` | `BackupTriggerStrategy` | String | 触发策略 |
| `cross_regions_enable` | `CrossRegionsEnable` | String | 跨地域备份开关（仅 Read 返回） |
| `cross_regions` | `CrossRegions` | List(String) | 跨地域目标地域（仅 Read 返回） |

### Read 模块

在 `resourceTencentCloudCynosdbBackupConfigRead` 中，读取 `DescribeBackupConfig` 返回的 `SnapshotSecondaryBackupConfig`（`BackupConfigInfo` 类型）并 Set 到 state。

### Update 模块

在 `resourceTencentCloudCynosdbBackupConfigUpdate` 中：
1. 将 `snapshot_secondary_backup_config` 加入 `mutableArgs` 列表
2. 读取 TF 配置中的 `snapshot_secondary_backup_config`，构建 `SnapshotBackupConfig` 对象
3. 设置到 `request.SnapshotSecondaryBackupConfig`

## Backward Compatibility

- ✅ 新增 Optional 字段，不影响现有配置
- ✅ 旧版 State 不含该字段，升级后首次 plan/apply 会将云端实际值写入 state
- ✅ 不修改任何已有字段

## Files to Modify

1. `tencentcloud/services/cynosdb/resource_tc_cynosdb_backup_config.go` — schema、read、update
2. `tencentcloud/services/cynosdb/resource_tc_cynosdb_backup_config.md` — 示例更新
