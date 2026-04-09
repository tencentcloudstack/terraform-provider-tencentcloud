# 变更提案：tencentcloud_tdmq_role 新增 token 字段

## 变更类型

**功能增强** — 在 `tencentcloud_tdmq_role` 资源中补充 `DescribeRoles` 接口返回的 `Token` 字段。

## Why

`DescribeRoles` 接口（https://cloud.tencent.com/document/api/1179/62399）在 `RoleSets` 中返回 `Token` 字段，该字段是角色的认证令牌，用于身份验证和授权。

当前 Read 模块只读取了 `RoleName` 和 `Remark`，未将 `Token` 写入 state，导致用户无法通过 Terraform 获取角色 token 用于后续配置（如消息队列客户端鉴权）。

## What Changes

| 文件 | 修改内容 |
|------|---------|
| `tencentcloud/services/tpulsar/resource_tc_tdmq_role.go` | schema 新增 `token`（Computed, Sensitive）；Read 模块新增 `info.Token` 赋值 |

### 字段设计

| Terraform 字段 | SDK 字段 | 类型 | 属性 | 说明 |
|--------------|---------|------|------|------|
| `token` | `Token` | `string` | Computed, Sensitive | 角色认证令牌，由接口返回，不可手动设置 |

- `Computed`：token 由服务端生成，用户不填写
- `Sensitive`：token 是鉴权凭据，需要在 plan/apply 输出中脱敏

### 向后兼容性

✅ 完全向后兼容：`token` 为纯 Computed 字段，不影响已有配置。
