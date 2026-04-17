# 变更提案：tencentcloud_redis_instance 支持 charge_type 字段修改

## 变更类型

**功能增强** — 将 `charge_type` 字段从 `ForceNew` 改为可更新，通过调用 `ModifyInstanceChargeType` 接口实现按量计费与包年包月之间的互转。

## Why

当前 `charge_type` 字段设置了 `ForceNew: true`，修改计费类型会触发实例删除重建，对于生产环境来说代价极高。

腾讯云 Redis 已支持 `ModifyInstanceChargeType` 接口（https://cloud.tencent.com/document/api/239/130332），可以在不重建实例的情况下修改计费类型。

## 接口入参

| 参数 | 必选 | 类型 | 说明 |
|------|------|------|------|
| `InstanceIds` | 是 | []String | 实例 ID 数组（本场景传单个 ID） |
| `InstanceChargeType` | 是 | String | `PREPAID`：按量转包年包月；`POSTPAID`：包年包月转按量 |
| `Period` | 否 | Integer | 购买时长（月），仅 `PREPAID` 时需要，范围 [1, 36] |

## SDK 状态

⚠️ **当前 vendor 目录中的 SDK 不包含 `ModifyInstanceChargeType` 方法**，需要通过直接构造 HTTP 请求或等待 SDK 更新。本变更采用**直接调用 SDK client 的 `SendOctetStream` 或使用 `NewCommonRequest`** 的方式，或者使用 SDK 的 `NewModifyInstanceChargeTypeRequest`（如 vendor 升级后）。

实际实现方案：在 service 层使用 `sdk.NewModifyInstanceChargeTypeRequest` 方式，如果 vendor 中不存在，则通过**升级 vendor** 引入。需先确认 vendor 升级路径，本变更假设通过运行 `go get` 更新 SDK vendor 来解决。

## What Changes

### 修改位置

| 文件 | 修改内容 |
|------|---------|
| `tencentcloud/services/crs/resource_tc_redis_instance.go` | `charge_type` 去掉 `ForceNew`，description 更新；update 函数新增 `charge_type` 变更处理 |
| `tencentcloud/services/crs/service_tencentcloud_redis.go` | 新增 `ModifyInstanceChargeType` service 方法 |

### charge_type 更新逻辑

- `POSTPAID → PREPAID`：需要 `prepaid_period`，调接口后等待实例 online
- `PREPAID → POSTPAID`：不需要 `Period`，调接口后等待实例 online
- `prepaid_period` 字段从 unsupportedUpdateFields 中联动处理（与 charge_type 一起变更时允许）

### 向后兼容性

⚠️ 有轻微破坏性：去掉 `ForceNew` 后，已有 state 中 charge_type 的变更不再触发重建，而是原地修改。对于已有用户是正向变化。
