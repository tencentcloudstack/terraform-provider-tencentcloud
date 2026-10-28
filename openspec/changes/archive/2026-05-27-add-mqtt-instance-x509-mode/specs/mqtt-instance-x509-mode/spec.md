## ADDED Requirements

### Requirement: X509Mode 参数支持

资源 `tencentcloud_mqtt_instance` 必须支持 `x509_mode` 参数（Optional, Computed, TypeString），可选值为 `TLS`（单向认证）、`mTLS`（双向认证）、`BYOC`（一机一证），通过 ModifyInstance 接口设置和更新，通过 DescribeInstance 接口读取。

#### Scenario: Create instance with x509_mode specified
- **WHEN** 用户在资源配置中指定 `x509_mode`
- **THEN** 系统先创建实例并等待 RUNNING 状态
- **AND** 然后通过 ModifyInstance 接口设置 `X509Mode` 字段

#### Scenario: Create instance without x509_mode
- **WHEN** 用户未指定 `x509_mode`
- **THEN** 不通过 ModifyInstance 设置 `X509Mode`（除非有其他需要 Modify 的参数）
- **AND** Read 时从 DescribeInstance 读取 API 返回的默认值

#### Scenario: Read x509_mode from DescribeInstance response
- **WHEN** 执行资源的 Read 操作
- **THEN** 从 DescribeInstance 响应中读取 `X509Mode` 字段
- **AND** 如果不为 nil，设置到 state

#### Scenario: Update x509_mode
- **WHEN** 用户修改 `x509_mode` 的值
- **THEN** 通过 ModifyInstance 接口更新 `X509Mode` 字段
- **AND** 通过 Read（DescribeInstance）验证更新成功

#### Scenario: Import existing instance with x509_mode
- **WHEN** 用户导入已有 MQTT 实例
- **THEN** 正确读取并填充 `x509_mode` 到 state
