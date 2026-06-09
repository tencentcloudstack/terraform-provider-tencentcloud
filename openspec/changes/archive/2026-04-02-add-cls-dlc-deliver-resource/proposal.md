# 变更提案：新增 tencentcloud_cls_dlc_deliver 资源

## 变更类型

**新功能** — 在 CLS（日志服务）云产品下新增 DLC 投递任务资源，支持 Create / Read / Delete 操作（无 Update，删除后重建）。

## Why

腾讯云 CLS 支持将日志主题中的数据投递到数据湖计算（DLC）服务，用于离线分析和数据归档。目前 Provider 中缺少对应的 Terraform 资源，用户无法通过 IaC 方式管理 DLC 投递任务的生命周期。

### 接口信息

| 操作 | 接口名 | 文档 |
|------|--------|------|
| 创建 | `CreateDlcDeliver` | https://cloud.tencent.com/document/api/614/125886 |
| 查询 | `DescribeDlcDelivers` | https://cloud.tencent.com/document/api/614/125884 |
| 删除 | `DeleteDlcDeliver` | https://cloud.tencent.com/document/api/614/125885 |

SDK 包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`

### 注意事项

- `CreateDlcDeliver` 返回 `TaskId`，资源唯一 ID 设计为 `topicId#taskId`（`#` 分隔）
- 无 Update 接口，字段变更需要触发 `ForceNew` 或先删后建

## What Changes

### 资源 Schema 主要字段

**Required：**
- `topic_id`（String）：日志主题 ID
- `name`（String）：任务名称，≤64 字符
- `deliver_type`（Int）：投递类型，0=批投递，1=实时投递
- `start_time`（Int）：投递时间范围开始时间（Unix 时间戳）
- `dlc_info`（List，MaxItems=1）：DLC 配置信息，包含：
  - `table_info`（List，MaxItems=1）：表信息（`data_directory`、`database_name`、`table_name`）
  - `field_infos`（List）：字段映射（`cls_field`、`dlc_field`、`dlc_field_type`、`fill_field`、`disable`）
  - `partition_infos`（List）：分区信息（`cls_field`、`dlc_field`、`dlc_field_type`）
  - `partition_extra`（List，MaxItems=1）：分区额外信息（`time_format`、`time_zone`）

**Optional：**
- `max_size`（Int）：批投递文件大小（MB），`deliver_type=0` 时必填，范围 5~256
- `interval`（Int）：批投递间隔（秒），`deliver_type=0` 时必填，范围 300~900
- `end_time`（Int）：投递结束时间，空=不限时
- `has_services_log`（Int）：是否开启服务日志，1=关闭，2=开启

**Computed（只读）：**
- `task_id`（String）：投递任务 ID

### 资源 ID

`d.SetId(topicId + "#" + taskId)`，Read/Delete 时通过 `strings.Split(d.Id(), "#")` 解析。

### 新增文件

| 文件 | 说明 |
|------|------|
| `tencentcloud/services/cls/resource_tc_cls_dlc_deliver.go` | 资源主文件（Schema + CRD） |
| `tencentcloud/services/cls/resource_tc_cls_dlc_deliver.md` | 资源文档示例 |
| `tencentcloud/services/cls/resource_tc_cls_dlc_deliver_test.go` | 验收测试 |

### 修改文件

| 文件 | 修改内容 |
|------|---------|
| `tencentcloud/services/cls/service_tencentcloud_cls.go` | 新增 `DescribeClsDlcDeliverById` service 方法 |
| `tencentcloud/provider.go` | 注册新资源 `tencentcloud_cls_dlc_deliver` |

### 向后兼容性

✅ 纯新增，不影响任何现有资源。

### 代码风格参考

严格参照 `tencentcloud_cls_scheduled_sql` 资源风格（CLS 包内已有代码风格），同时参考 `tencentcloud_igtm_strategy` 的 `resource.Retry` 写法：
- `defer tccommon.LogElapsed(...)()` + `defer tccommon.InconsistentCheck(d, meta)()`
- `resource.Retry(tccommon.WriteRetryTimeout, ...)` / `resource.Retry(tccommon.ReadRetryTimeout, ...)`
- service 层封装 `DescribeClsDlcDeliverById`
