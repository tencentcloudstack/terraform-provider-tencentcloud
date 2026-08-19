## Context

`tencentcloud_mqtt_instance` 是管理 MQTT 实例的 Terraform 资源，支持 CRUD 操作。云 API `ModifyInstance` 接口支持 `X509Mode` 参数用于设置 TLS 认证模式，`DescribeInstance` 响应中返回 `X509Mode` 字段。

SDK 中相关定义：
- `ModifyInstanceRequestParams.X509Mode *string` — TLS：单向认证 / mTLS：双向认证 / BYOC：一机一证
- `DescribeInstanceResponseParams.X509Mode *string` — 同上

注意：`CreateInstance` 接口不支持 `X509Mode` 参数，因此需要在创建完成后通过 `ModifyInstance` 设置。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_mqtt_instance` 资源中新增 `x509_mode` 参数（Optional, Computed, TypeString）
- 在 Create 中，等待实例 RUNNING 后通过 ModifyInstance 设置 `X509Mode`
- 在 Update 中，当 `x509_mode` 发生变更时通过 ModifyInstance 更新
- 在 Read 中，从 DescribeInstance 响应读取并设置到 state
- 更新单元测试

**Non-Goals:**
- 不支持在 CreateInstance 接口中直接设置（API 不支持）
- 不修改其他 MQTT 相关资源

## Decisions

1. **Schema 定义**: `x509_mode` 为 `TypeString`、`Optional: true`、`Computed: true`。Optional 因为用户可以选择性配置；Computed 因为可从 API 读回。

2. **Create 实现**: 复用已有的 "创建后 ModifyInstance" 模式。在实例创建并等待 RUNNING 状态后，与 `automatic_activation`/`authorization_policy` 一起通过 ModifyInstance 设置。

3. **Update 实现**: 将 `x509_mode` 加入 mutableArgs 列表，在 ModifyInstance 请求中设置 `request.X509Mode`。

4. **Read 实现**: 当 `respData.X509Mode` 不为 nil 时，设置到 state。

## Risks / Trade-offs

- [风险] 创建后立即修改可能在极端情况下失败 → 已有重试机制覆盖
- [风险] `X509Mode` 设置为 `mTLS` 或 `BYOC` 时可能需要先配置证书 → 由用户确保前置条件满足
