## Why

MQTT 实例的 ModifyInstance 接口和 DescribeInstance 接口支持 `X509Mode` 字段（TLS 认证模式：TLS 单向认证 / mTLS 双向认证 / BYOC 一机一证），需要在 Terraform 资源 `tencentcloud_mqtt_instance` 中支持该参数，使用户可以通过 Terraform 配置和修改实例的 X509 认证模式。

## What Changes

- 在 `tencentcloud_mqtt_instance` 资源的 Schema 中新增 `x509_mode` 参数（Optional, Computed, TypeString），可选值为 `TLS`、`mTLS`、`BYOC`
- 在资源的 Create 方法中，创建实例并等待 RUNNING 状态后，通过 ModifyInstance 接口设置 `x509_mode`
- 在资源的 Update 方法中，当 `x509_mode` 发生变更时，通过 ModifyInstance 接口更新
- 在资源的 Read 方法中，从 DescribeInstance 响应中读取 `X509Mode` 并设置到 state

## Capabilities

### New Capabilities

- `mqtt-instance-x509-mode`: 为 tencentcloud_mqtt_instance 资源新增 x509_mode 参数，支持通过 ModifyInstance 设置和修改 X509 认证模式，通过 DescribeInstance 读取当前值

### Modified Capabilities

## Impact

- `tencentcloud/services/mqtt/resource_tc_mqtt_instance.go`: 新增 schema 字段，修改 Create/Read/Update 逻辑
- `tencentcloud/services/mqtt/resource_tc_mqtt_instance_test.go`: 新增单元测试
- `.changelog/4136.txt`: 更新 changelog
